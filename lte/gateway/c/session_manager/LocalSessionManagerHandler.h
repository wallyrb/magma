/**
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
 */
#pragma once

#include <functional>

#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "LocalEnforcer.h"
#include "SessionID.h"
#include "SessionReporter.h"
#include "SessionStore.h"

using grpc::Server;
using grpc::ServerContext;
using grpc::Status;

namespace magma {
using namespace orc8r;

class LocalSessionManagerHandler {
 public:
  virtual ~LocalSessionManagerHandler() {}

  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  virtual void ReportRuleStats(
      ServerContext* context, const RuleRecordTable* request,
      std::function<void(Status, Void)> response_callback) = 0;

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  virtual void CreateSession(
      ServerContext* context, const LocalCreateSessionRequest* request,
      std::function<void(Status, LocalCreateSessionResponse)>
          response_callback) = 0;

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  virtual void EndSession(
      ServerContext* context, const LocalEndSessionRequest* request,
      std::function<void(Status, LocalEndSessionResponse)>
          response_callback) = 0;

  /**
   * Bind the returned bearer id to the policy for which it is created
   */
  virtual void BindPolicy2Bearer(
      ServerContext* context, const PolicyBearerBindingRequest* request,
      std::function<void(Status, PolicyBearerBindingResponse)>
          response_callback) = 0;
  /**
   * Update active rules for session
   */
  virtual void SetSessionRules(
      ServerContext* context, const SessionRules* request,
      std::function<void(Status, Void)> response_callback) = 0;
};

/**
 * LocalSessionManagerHandler processes proxied gRPC requests to the session
 * manager. The handler uses a monitor and reporter to keep track of state
 * and report to the cloud, respectively
 */
class LocalSessionManagerHandlerImpl : public LocalSessionManagerHandler {
 public:
  LocalSessionManagerHandlerImpl(
      std::shared_ptr<LocalEnforcer> monitor, SessionReporter* reporter,
      std::shared_ptr<AsyncDirectorydClient> directoryd_client,
      std::shared_ptr<EventsReporter> events_reporter,
      SessionStore& session_store);
  ~LocalSessionManagerHandlerImpl() {}
  /**
   * Report flow stats from pipelined and track the usage per rule
   */
  void ReportRuleStats(
      ServerContext* context, const RuleRecordTable* request,
      std::function<void(Status, Void)> response_callback);

  /**
   * Create a new session, initializing credit monitoring and requesting credit
   * from the cloud
   */
  void CreateSession(
      ServerContext* context, const LocalCreateSessionRequest* request,
      std::function<void(Status, LocalCreateSessionResponse)>
          response_callback);

  /**
   * Terminate a session, untracking credit and terminating in the cloud
   */
  void EndSession(
      ServerContext* context, const LocalEndSessionRequest* request,
      std::function<void(Status, LocalEndSessionResponse)> response_callback);

  /**
   * Bind the returned bearer id to the policy for which it is created; if
   * the returned bearer id is 0 then the dedicated bearer request is rejected
   */
  void BindPolicy2Bearer(
      ServerContext* context, const PolicyBearerBindingRequest* request,
      std::function<void(Status, PolicyBearerBindingResponse)>
          response_callback);

  /**
   * Update active rules for session
   * Get the SessionMap for the updates, apply the set rules and update the
   * store. The rule updates should be also propagated to PipelineD
   */
  void SetSessionRules(
      ServerContext* context, const SessionRules* request,
      std::function<void(Status, Void)> response_callback);

 private:
  SessionStore& session_store_;
  std::shared_ptr<LocalEnforcer> enforcer_;
  SessionReporter* reporter_;
  std::shared_ptr<AsyncDirectorydClient> directoryd_client_;
  std::shared_ptr<EventsReporter> events_reporter_;
  SessionIDGenerator id_gen_;
  uint64_t current_epoch_;
  uint64_t reported_epoch_;
  std::chrono::milliseconds retry_timeout_;
  static const std::string hex_digit_;

 private:
  void check_usage_for_reporting(
      SessionMap session_map, SessionUpdate& session_uc);
  bool is_pipelined_restarted();
  bool restart_pipelined(const std::uint64_t& epoch);

  void end_session(
      SessionMap& session_map, const SubscriberID& sid, const std::string& apn,
      std::function<void(Status, LocalEndSessionResponse)> response_callback);

  std::string convert_mac_addr_to_str(const std::string& mac_addr);

  void add_session_to_directory_record(
      const std::string& imsi, const std::string& session_id,
      const std::string& msisdn);

  /**
   * handle_create_session_cwf handles a sequence of actions needed for the
   * RATType=WLAN case. It is responsible for responding to the original
   * LocalCreateSession request.
   * If there is an existing session for the IMSI that is ACTIVE, we will
   * simply update its SessionConfig with the new context. In this case, we will
   * NOT send a CreateSession request into FeG/PolicyDB.
   * Otherwise, we will go through the procedure of creating a new context.
   * @param session_map - SessionMap that contains all sessions with IMSI
   * @param sid - newly created SessionID
   * @param cfg - newly created SessionConfig from the LocalCreateSessionRequest
   * @param cb - callback needed to respond to the original
   * LocalCreateSessionRequest
   */
  void handle_create_session_cwf(
      SessionMap& session_map, const std::string& sid, SessionConfig cfg,
      std::function<void(Status, LocalCreateSessionResponse)> cb);

  /**
   * Handle the logic to recycle an existing CWF session. This involves updating
   * the existing SessionConfig with the new one. This function is responsible
   * for responding to the original LocalCreateSession call with the cb.
   */
  void recycle_cwf_session(
      const std::string& imsi, const std::string& sid, const SessionConfig& cfg,
      SessionMap& session_map,
      std::function<void(Status, LocalCreateSessionResponse)> cb);
  /**
   * handle_create_session_lte handles a sequence of actions needed for the
   * RATType=LTE case. It is responsible for responding to the original
   * LocalCreateSession request.
   * If there is an existing identical session, same SessionConfig, for the IMSI
   * that is ACTIVE, we will reuse this session. In this case, we will NOT send
   * a CreateSession request into FeG/PolicyDB.
   * Otherwise, we will go through the procedure of creating a new context.
   * @param session_map - SessionMap that contains all sessions with IMSI
   * @param sid - newly created SessionID
   * @param cfg - newly created SessionConfig from the LocalCreateSessionRequest
   * @param cb - callback needed to respond to the original
   * LocalCreateSessionRequest
   */
  void handle_create_session_lte(
      SessionMap& session_map, const std::string& sid, SessionConfig cfg,
      std::function<void(Status, LocalCreateSessionResponse)> cb);

  /**
   * Send session creation request to the CentralSessionController.
   * If it is successful, create a session in session_map, and respond to
   * gRPC caller.
   */
  void send_create_session(
      SessionMap& session_map, const std::string& sid, const SessionConfig& cfg,
      std::function<void(grpc::Status, LocalCreateSessionResponse)> cb);

  void handle_setup_callback(
      const std::uint64_t& epoch, Status status, SetupFlowsResult resp);

  void report_session_update_event(
      SessionMap& session_map, SessionUpdate& session_update);

  void report_session_update_event_failure(
      SessionMap& session_map, SessionUpdate& session_update,
      const std::string& failure_reason);

  void send_local_create_session_response(
      Status status, const std::string& sid,
      std::function<void(Status, LocalCreateSessionResponse)>
          response_callback);

  void log_create_session(SessionConfig& cfg);

  std::string bytes_to_hex(const std::string& s);
};

}  // namespace magma
