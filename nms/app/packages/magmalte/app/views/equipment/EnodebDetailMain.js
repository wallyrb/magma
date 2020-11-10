/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebConfig from './EnodebDetailConfig';
import EnodebContext from '../../components/context/EnodebContext';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {EnodebJsonConfig} from './EnodebDetailConfig';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {Redirect, Route, Switch} from 'react-router-dom';
import {RunGatewayCommands} from '../../state/lte/EquipmentState';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
}));
const CHART_TITLE = 'Bandwidth Usage';

export function EnodebDetail() {
  const ctx = useContext(EnodebContext);
  const {relativePath, relativeUrl, match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];

  return (
    <>
      <TopBar
        header={`Equipment/${enbInfo.enb.name}`}
        tabs={[
          {
            label: 'Overview',
            to: '/overview',
            icon: DashboardIcon,
            filters: <EnodebRebootButton />,
          },
          {
            label: 'Config',
            to: '/config',
            icon: SettingsIcon,
            filters: <EnodebRebootButton />,
          },
        ]}
      />

      <Switch>
        <Route path={relativePath('/overview')} component={Overview} />
        <Route
          path={relativePath('/config/json')}
          component={EnodebJsonConfig}
        />
        <Route path={relativePath('/config')} component={EnodebConfig} />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function EnodebRebootButtonInternal(props: WithAlert) {
  const classes = useStyles();
  const ctx = useContext(EnodebContext);
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const gatewayId = enbInfo?.enb_state?.reporting_gateway_id;
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleClick = () => {
    if (gatewayId == null) {
      enqueueSnackbar('Unable to trigger reboot, reporting gateway not found', {
        variant: 'error',
      });
      return;
    }

    props
      .confirm(`Are you sure you want to reboot ${enodebSerial}?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }
        const params = {
          command: 'reboot_enodeb',
          params: {shell_params: {[enodebSerial]: {}}},
        };

        try {
          await RunGatewayCommands({
            networkId,
            gatewayId,
            command: 'generic',
            params,
          });
          enqueueSnackbar('eNodeB reboot triggered successfully', {
            variant: 'success',
          });
        } catch (e) {
          enqueueSnackbar(e.response?.data?.message ?? e.message, {
            variant: 'error',
          });
        }
      });
  };

  return (
    <Button
      variant="contained"
      className={classes.appBarBtn}
      onClick={handleClick}>
      Reboot
    </Button>
  );
}
const EnodebRebootButton = withAlert(EnodebRebootButtonInternal);

function Overview() {
  const ctx = useContext(EnodebContext);
  const classes = useStyles();
  const {match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  const enbInfo = ctx.state.enbInfo[enodebSerial];
  const enbIpAddress = enbInfo?.enb_state?.ip_address;
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label={enbInfo.enb.name}
              />
              <EnodebSummary />
            </Grid>

            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <EnodebStatus />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <DateTimeMetricChart
            title={CHART_TITLE}
            unit={'Throughput(mb/s)'}
            queries={[
              `rate(gtp_port_user_plane_dl_bytes{service="pipelined", ip_addr="${enbIpAddress}"}[5m])/1000`,
              `rate(gtp_port_user_plane_ul_bytes{service="pipelined", ip_addr="${enbIpAddress}"}[5m])/1000`,
            ]}
            legendLabels={['Download', 'Upload']}
          />
        </Grid>
      </Grid>
    </div>
  );
}

export default EnodebDetail;
