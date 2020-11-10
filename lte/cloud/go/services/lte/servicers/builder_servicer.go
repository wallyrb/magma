/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package servicers

import (
	"context"
	"fmt"
	"sort"

	"magma/lte/cloud/go/lte"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/thoas/go-funk"
)

type builderServicer struct{}

func NewBuilderServicer() builder_protos.MconfigBuilderServer {
	return &builderServicer{}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (*builder_protos.BuildResponse, error) {
	ret := &builder_protos.BuildResponse{ConfigsByKey: map[string][]byte{}}

	network, err := (configurator.Network{}).FromStorageProto(request.Network)
	if err != nil {
		return nil, err
	}
	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.Graph)
	if err != nil {
		return nil, err
	}

	// Only build an mconfig if cellular network and gateway configs exist
	inwConfig, found := network.Configs[lte.CellularNetworkConfigType]
	if !found || inwConfig == nil {
		return ret, nil
	}
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)

	cellGW, err := graph.GetEntity(lte.CellularGatewayEntityType, request.GatewayId)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	if cellGW.Config == nil {
		return ret, nil
	}
	cellularGwConfig := cellGW.Config.(*lte_models.GatewayCellularConfigs)

	if err := validateConfigs(cellularNwConfig, cellularGwConfig); err != nil {
		return nil, err
	}

	enodebs, err := graph.GetAllChildrenOfType(cellGW, lte.CellularEnodebEntityType)
	if err != nil {
		return nil, err
	}

	gwRan := cellularGwConfig.Ran
	gwEpc := cellularGwConfig.Epc
	gwNonEpsService := cellularGwConfig.NonEpsService
	nwRan := cellularNwConfig.Ran
	nwEpc := cellularNwConfig.Epc
	nonEPSServiceMconfig := getNonEPSServiceMconfigFields(gwNonEpsService)

	pipelineDServices, err := getPipelineDServicesConfig(nwEpc.NetworkServices)
	if err != nil {
		return nil, err
	}

	enbConfigsBySerial := getEnodebConfigsBySerial(cellularNwConfig, cellularGwConfig, enodebs)

	vals := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel:            protos.LogLevel_INFO,
			Pci:                 int32(gwRan.Pci),
			FddConfig:           getFddConfig(nwRan.FddConfig),
			TddConfig:           getTddConfig(nwRan.TddConfig),
			BandwidthMhz:        int32(nwRan.BandwidthMhz),
			AllowEnodebTransmit: swag.BoolValue(gwRan.TransmitEnabled),
			Tac:                 int32(nwEpc.Tac),
			PlmnidList:          fmt.Sprintf("%s%s", nwEpc.Mcc, nwEpc.Mnc),
			CsfbRat:             nonEPSServiceMconfig.csfbRat,
			Arfcn_2G:            nonEPSServiceMconfig.arfcn_2g,
			EnbConfigsBySerial:  enbConfigsBySerial,
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel:        protos.LogLevel_INFO,
			IpBlock:         gwEpc.IPBlock,
			IpAllocatorType: getMobilityDIPAllocator(nwEpc),
			StaticIpEnabled: getMobilityDStaticIPAllocation(nwEpc),
			MultiApnIpAlloc: getMobilityDMultuAPNIPAlloc(nwEpc),
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      nwEpc.Mcc,
			Mnc:                      nwEpc.Mnc,
			Tac:                      int32(nwEpc.Tac),
			MmeCode:                  1,
			MmeGid:                   1,
			EnableDnsCaching:         shouldEnableDNSCaching(cellularGwConfig.DNS),
			NonEpsServiceControl:     nonEPSServiceMconfig.nonEpsServiceControl,
			CsfbMcc:                  nonEPSServiceMconfig.csfbMcc,
			CsfbMnc:                  nonEPSServiceMconfig.csfbMnc,
			Lac:                      nonEPSServiceMconfig.lac,
			HssRelayEnabled:          swag.BoolValue(nwEpc.HssRelayEnabled),
			CloudSubscriberdbEnabled: nwEpc.CloudSubscriberdbEnabled,
			AttachedEnodebTacs:       getEnodebTacs(enbConfigsBySerial),
			DnsPrimary:               gwEpc.DNSPrimary,
			DnsSecondary:             gwEpc.DNSSecondary,
			NatEnabled:               swag.BoolValue(gwEpc.NatEnabled),
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:                 protos.LogLevel_INFO,
			UeIpBlock:                gwEpc.IPBlock,
			NatEnabled:               swag.BoolValue(gwEpc.NatEnabled),
			DefaultRuleId:            nwEpc.DefaultRuleID,
			Services:                 pipelineDServices,
			SgiManagementIfaceVlan:   gwEpc.SgiManagementIfaceVlan,
			SgiManagementIfaceIpAddr: gwEpc.SgiManagementIfaceStaticIP,
			SgiManagementIfaceGw:     gwEpc.SgiManagementIfaceGw,
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:        protos.LogLevel_INFO,
			LteAuthOp:       nwEpc.LteAuthOp,
			LteAuthAmf:      nwEpc.LteAuthAmf,
			SubProfiles:     getSubProfiles(nwEpc),
			HssRelayEnabled: swag.BoolValue(nwEpc.HssRelayEnabled),
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: swag.BoolValue(nwEpc.GxGyRelayEnabled),
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
		},
		"dnsd": getGatewayCellularDNSMConfig(cellularGwConfig.DNS),
	}

	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func validateConfigs(nwConfig *lte_models.NetworkCellularConfigs, gwConfig *lte_models.GatewayCellularConfigs) error {
	if nwConfig == nil {
		return errors.New("Cellular network config is nil")
	}
	if gwConfig == nil {
		return errors.New("Cellular gateway config is nil")
	}

	if gwConfig.Ran == nil {
		return errors.New("Gateway RAN config is nil")
	}
	if gwConfig.Epc == nil {
		return errors.New("Gateway EPC config is nil")
	}
	if nwConfig.Ran == nil {
		return errors.New("Network RAN config is nil")
	}
	if nwConfig.Epc == nil {
		return errors.New("Network EPC config is nil")
	}
	return nil
}

type nonEPSServiceMconfigFields struct {
	csfbRat              lte_mconfig.EnodebD_CSFBRat
	arfcn_2g             []int32
	nonEpsServiceControl lte_mconfig.MME_NonEPSServiceControl
	csfbMcc              string
	csfbMnc              string
	lac                  int32
}

func getNonEPSServiceMconfigFields(gwNonEpsService *lte_models.GatewayNonEpsConfigs) nonEPSServiceMconfigFields {
	if gwNonEpsService == nil {
		return nonEPSServiceMconfigFields{
			csfbRat:              lte_mconfig.EnodebD_CSFBRAT_2G,
			arfcn_2g:             []int32{},
			nonEpsServiceControl: lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			csfbMcc:              "",
			csfbMnc:              "",
			lac:                  1,
		}
	} else {
		arfcn2g := make([]int32, 0, len(gwNonEpsService.Arfcn2g))
		for _, i := range gwNonEpsService.Arfcn2g {
			arfcn2g = append(arfcn2g, int32(i))
		}

		return nonEPSServiceMconfigFields{
			csfbRat:              lte_mconfig.EnodebD_CSFBRat(swag.Uint32Value(gwNonEpsService.CsfbRat)),
			arfcn_2g:             arfcn2g,
			nonEpsServiceControl: lte_mconfig.MME_NonEPSServiceControl(swag.Uint32Value(gwNonEpsService.NonEpsServiceControl)),
			csfbMcc:              gwNonEpsService.CsfbMcc,
			csfbMnc:              gwNonEpsService.CsfbMnc,
			lac:                  int32(swag.Uint32Value(gwNonEpsService.Lac)),
		}
	}
}

var networkServicesByName = map[string]lte_mconfig.PipelineD_NetworkServices{
	"metering":           lte_mconfig.PipelineD_METERING,
	"dpi":                lte_mconfig.PipelineD_DPI,
	"policy_enforcement": lte_mconfig.PipelineD_ENFORCEMENT,
}

// move this out of this package eventually
func getPipelineDServicesConfig(networkServices []string) ([]lte_mconfig.PipelineD_NetworkServices, error) {
	if networkServices == nil || len(networkServices) == 0 {
		return []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_ENFORCEMENT,
		}, nil
	}
	apps := make([]lte_mconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[service]
		if !found {
			return nil, errors.Errorf("unknown network service name %s", service)
		}
		apps = append(apps, mc)
	}
	return apps, nil
}

func getFddConfig(fddConfig *lte_models.NetworkRanConfigsFddConfig) *lte_mconfig.EnodebD_FDDConfig {
	if fddConfig == nil {
		return nil
	}
	return &lte_mconfig.EnodebD_FDDConfig{
		Earfcndl: int32(fddConfig.Earfcndl),
		Earfcnul: int32(fddConfig.Earfcnul),
	}
}

func getTddConfig(tddConfig *lte_models.NetworkRanConfigsTddConfig) *lte_mconfig.EnodebD_TDDConfig {
	if tddConfig == nil {
		return nil
	}

	return &lte_mconfig.EnodebD_TDDConfig{
		Earfcndl:               int32(tddConfig.Earfcndl),
		SubframeAssignment:     int32(tddConfig.SubframeAssignment),
		SpecialSubframePattern: int32(tddConfig.SpecialSubframePattern),
	}
}

func getEnodebConfigsBySerial(nwConfig *lte_models.NetworkCellularConfigs, gwConfig *lte_models.GatewayCellularConfigs, enodebs []configurator.NetworkEntity) map[string]*lte_mconfig.EnodebD_EnodebConfig {
	ret := make(map[string]*lte_mconfig.EnodebD_EnodebConfig, len(enodebs))
	for _, ent := range enodebs {
		serial := ent.Key
		ienbConfig := ent.Config
		if ienbConfig == nil {
			glog.Errorf("enb with serial %s is missing config", serial)
		}

		enodebConfig := ienbConfig.(*lte_models.EnodebConfig)
		enbMconfig := &lte_mconfig.EnodebD_EnodebConfig{}

		if enodebConfig.ConfigType == "MANAGED" {
			cellularEnbConfig := enodebConfig.ManagedConfig
			enbMconfig.Earfcndl = int32(cellularEnbConfig.Earfcndl)
			enbMconfig.SubframeAssignment = int32(cellularEnbConfig.SubframeAssignment)
			enbMconfig.SpecialSubframePattern = int32(cellularEnbConfig.SpecialSubframePattern)
			enbMconfig.Pci = int32(cellularEnbConfig.Pci)
			enbMconfig.TransmitEnabled = swag.BoolValue(cellularEnbConfig.TransmitEnabled)
			enbMconfig.DeviceClass = cellularEnbConfig.DeviceClass
			enbMconfig.BandwidthMhz = int32(cellularEnbConfig.BandwidthMhz)
			enbMconfig.Tac = int32(cellularEnbConfig.Tac)
			enbMconfig.CellId = int32(swag.Uint32Value(cellularEnbConfig.CellID))

			// override zero values with network/gateway configs
			if enbMconfig.Earfcndl == 0 {
				enbMconfig.Earfcndl = int32(nwConfig.GetEarfcndl())
			}
			if enbMconfig.SubframeAssignment == 0 {
				if nwConfig.Ran.TddConfig != nil {
					enbMconfig.SubframeAssignment = int32(nwConfig.Ran.TddConfig.SubframeAssignment)
				}
			}
			if enbMconfig.SpecialSubframePattern == 0 {
				if nwConfig.Ran.TddConfig != nil {
					enbMconfig.SpecialSubframePattern = int32(nwConfig.Ran.TddConfig.SpecialSubframePattern)
				}
			}
			if enbMconfig.Pci == 0 {
				enbMconfig.Pci = int32(gwConfig.Ran.Pci)
			}
			if enbMconfig.BandwidthMhz == 0 {
				enbMconfig.BandwidthMhz = int32(nwConfig.Ran.BandwidthMhz)
			}
			if enbMconfig.Tac == 0 {
				enbMconfig.Tac = int32(nwConfig.Epc.Tac)
			}

		} else if enodebConfig.ConfigType == "UNMANAGED" {
			cellularEnbConfig := enodebConfig.UnmanagedConfig
			enbMconfig.CellId = int32(swag.Uint32Value(cellularEnbConfig.CellID))
			enbMconfig.Tac = int32(swag.Uint32Value(cellularEnbConfig.Tac))
			enbMconfig.IpAddress = string(*cellularEnbConfig.IPAddress)

			if enbMconfig.Tac == 0 {
				enbMconfig.Tac = int32(nwConfig.Epc.Tac)
			}
		}

		ret[serial] = enbMconfig
	}
	return ret
}

func getEnodebTacs(enbConfigsBySerial map[string]*lte_mconfig.EnodebD_EnodebConfig) []int32 {
	ret := make([]int32, 0, len(enbConfigsBySerial))
	for _, enbConfig := range enbConfigsBySerial {
		ret = append(ret, enbConfig.Tac)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i] < ret[j] })
	return ret
}

func getSubProfiles(epc *lte_models.NetworkEpcConfigs) map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile {
	if epc.SubProfiles == nil {
		return map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile{}
	}

	ret := map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile{}
	for name, profile := range epc.SubProfiles {
		ret[name] = &lte_mconfig.SubscriberDB_SubscriptionProfile{
			MaxUlBitRate: profile.MaxUlBitRate,
			MaxDlBitRate: profile.MaxDlBitRate,
		}
	}
	return ret
}

func getMobilityDIPAllocator(epc *lte_models.NetworkEpcConfigs) lte_mconfig.MobilityD_IpAllocatorType {
	if epc.Mobility == nil {
		return lte_mconfig.MobilityD_IP_POOL
	}
	if epc.Mobility.IPAllocationMode == lte_models.DHCPBroadcastAllocationMode {
		return lte_mconfig.MobilityD_DHCP
	}
	// For other modes set IP pool allocator
	return lte_mconfig.MobilityD_IP_POOL
}

func getMobilityDStaticIPAllocation(epc *lte_models.NetworkEpcConfigs) bool {
	if epc.Mobility == nil {
		return false
	}
	return epc.Mobility.EnableStaticIPAssignments
}

func getMobilityDMultuAPNIPAlloc(epc *lte_models.NetworkEpcConfigs) bool {
	if epc.Mobility == nil {
		return false
	}
	return epc.Mobility.EnableMultiApnIPAllocation
}

func getGatewayCellularDNSMConfig(gwDns *lte_models.GatewayDNSConfigs) *lte_mconfig.DnsD {
	if gwDns == nil {
		return &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
			EnableCaching:     false,
			LocalTTL:          0,
			Records:           []*lte_mconfig.GatewayDNSConfigRecordsItems{},
		}
	} else {
		return &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: swag.BoolValue(gwDns.DhcpServerEnabled),
			EnableCaching:     shouldEnableDNSCaching(gwDns),
			LocalTTL:          *gwDns.LocalTTL,
			Records:           getGatewayDnsRecords(gwDns),
		}
	}
}

func getGatewayDnsRecords(dns *lte_models.GatewayDNSConfigs) []*lte_mconfig.GatewayDNSConfigRecordsItems {
	if dns.Records == nil {
		return []*lte_mconfig.GatewayDNSConfigRecordsItems{}
	}

	ret := make([]*lte_mconfig.GatewayDNSConfigRecordsItems, 0, len(dns.Records))
	for _, record := range dns.Records {
		recordProto := &lte_mconfig.GatewayDNSConfigRecordsItems{}
		recordProto.Domain = record.Domain
		recordProto.ARecord = funk.Map(record.ARecord, func(a strfmt.IPv4) string { return string(a) }).([]string)
		recordProto.AaaaRecord = funk.Map(record.AaaaRecord, func(a strfmt.IPv6) string { return string(a) }).([]string)
		recordProto.CnameRecord = make([]string, 0, len(record.CnameRecord))
		for _, cRecord := range record.CnameRecord {
			recordProto.CnameRecord = append(recordProto.CnameRecord, cRecord)
		}
		ret = append(ret, recordProto)
	}
	return ret
}

func shouldEnableDNSCaching(dns *lte_models.GatewayDNSConfigs) bool {
	if dns == nil {
		return false
	}
	return swag.BoolValue(dns.EnableCaching)
}
