# Phase 4: Multiplayer Physics Sync — 작업지시서

_작성일: 2026-04-18_
_속한 로드맵: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)_
_예상 소요: 1-2주_

## 페이즈 목표

2-4명의 사용자가 같은 room에서 동시에 조립하고 테스트해도 서버 상태가 깨지지 않도록 동시성, 충돌 처리, snapshot 효율을 강화한다. 이 단계는 “멀티플레이어에서 견딜 수 있는 구조인가”를 검증하는 안정화 단계다.

## 전제 조건

- [x] Phase 3 완료 및 사용자 승인
- [x] 서버 snapshot과 command protocol이 기본 동작 중
- [x] actuator transform replication이 구현됨

## 포함 기능

1. room/session model
2. command queue와 conflict resolution
3. snapshot delta 또는 changed-set 전송
4. basic collision/constraint policy
5. tick performance metrics
6. multiplayer manual test scenario

## 이 페이즈에서 하지 않는 것

- public lobby/matchmaking → Later
- 권한/밴/신고 시스템 → Out of Scope
- 대규모 interest management → Later

## 작업 체크리스트

### 작업 그룹 A: Room/session model

- [x] **T4.A.1** — room model 구현.
  - 파일:
    - `internal/server/room.go`
    - `internal/server/room_test.go`
  - 내용:
    - `RoomID`
    - room별 simulation instance
    - client join/leave
    - 기본 room 하나로 시작하되 다중 room 확장 가능하게 분리
  - 검증: `go test ./internal/server/...`

- [x] **T4.A.2** — client identity와 presence snapshot 추가.
  - 파일:
    - `internal/server/client.go`
    - `internal/netproto/snapshot.go`
    - `web/src/net/protocol.ts`
  - 내용:
    - 임시 nickname 또는 client index
    - connected clients count
    - debug overlay 표시
  - 검증: `go test ./... && cd web && npm run build`

### 작업 그룹 B: Command conflict handling

- [x] **T4.B.1** — command ordering 정책 문서화 및 구현.
  - 파일:
    - `internal/sim/commands.go`
    - `internal/sim/commands_test.go`
  - 내용:
    - tick 내 command sort 기준: received sequence, client id, command id
    - 같은 좌표 edit conflict 처리
    - duplicate command id 무시
  - 테스트:
    - 동일 좌표 동시 place
    - place/remove 충돌
    - duplicate command 재전송
  - 검증: `go test ./internal/sim/...`

- [ ] **T4.B.2** — client command acknowledgement.
  - 파일:
    - `internal/netproto/snapshot.go`
    - `web/src/net/socket.ts`
  - 내용:
    - snapshot에 accepted/rejected command ids 또는 last processed command 포함
    - rejected reason은 debug overlay에 표시
  - 검증: 브라우저에서 충돌 command의 reject 표시 확인

### 작업 그룹 C: Snapshot efficiency

- [ ] **T4.C.1** — changed-set snapshot 구현.
  - 파일:
    - `internal/netproto/snapshot.go`
    - `internal/sim/snapshot_builder.go`
    - `internal/sim/snapshot_builder_test.go`
  - 내용:
    - full snapshot과 delta snapshot 구분
    - changed blocks
    - removed blocks
    - changed entities
    - periodic full snapshot fallback
  - 검증: `go test ./internal/sim/...`

- [ ] **T4.C.2** — client delta apply 구현.
  - 파일:
    - `web/src/state/worldStore.ts`
    - `web/src/state/snapshotStore.ts`
  - 내용:
    - full snapshot reset
    - delta block/entity update
    - missing base snapshot 감지 시 resync 요청 또는 full snapshot 대기
  - 검증: `cd web && npm test`

### 작업 그룹 D: Basic physics constraints

- [ ] **T4.D.1** — 단순 충돌/점유 규칙 구현.
  - 파일:
    - `internal/physics/collision.go`
    - `internal/physics/collision_test.go`
  - 내용:
    - solid block 점유 좌표에 actuator head가 진입하지 못하게 제한
    - collision response는 stop 또는 clamp로 단순화
  - 검증: `go test ./internal/physics/...`

- [ ] **T4.D.2** — actuator와 collision policy 통합.
  - 파일:
    - `internal/actuator/piston.go`
    - `internal/sim/simulation.go`
  - 내용:
    - 피스톤 이동 target이 막히면 정지
    - debug snapshot에 blocked 상태 포함
  - 검증: `go test ./...`

### 작업 그룹 E: Observability

- [ ] **T4.E.1** — 서버 metrics 로그 추가.
  - 파일:
    - `internal/sim/metrics.go`
    - `internal/server/server.go`
  - 내용:
    - tick duration
    - command queue length
    - snapshot byte size
    - client count
  - 검증: 서버 로그에서 주기적 metrics 확인

- [ ] **T4.E.2** — multiplayer smoke scenario 작성.
  - 파일:
    - `docs/manual-tests/multiplayer-smoke.md`
  - 내용:
    - 클라이언트 4개 연결
    - 동시 block edit
    - 버튼-피스톤 demo
    - snapshot delay 관찰
  - 검증: 문서 절차대로 수동 확인

---

## ✅ Phase 4 Checkpoint

**구현 확인:**
- [ ] room/session 구조가 분리되어 있다.
- [x] command conflict가 결정적으로 처리된다.
- [ ] delta snapshot 또는 changed-set snapshot이 동작한다.
- [ ] basic collision/constraint로 actuator가 solid block을 뚫지 않는다.
- [ ] tick/snapshot/client metrics를 확인할 수 있다.

**자동 검증:**
- [ ] 서버 테스트 통과: `go test ./...`
- [ ] 클라이언트 테스트 통과: `cd web && npm test`
- [ ] 클라이언트 빌드 통과: `cd web && npm run build`

**수동 확인:**
- [ ] 4개 브라우저 창이 같은 room에 접속된다.
- [ ] 2명이 같은 좌표를 동시에 수정해도 모든 클라이언트가 같은 결과를 본다.
- [ ] 피스톤 앞에 solid block이 있으면 피스톤이 멈추거나 clamp된다.
- [ ] snapshot byte size와 tick duration이 debug/log로 확인된다.

**완료 처리:**
1. 위 항목 모두 통과 시 완료 요약과 검증 결과를 보고한다.
2. 사용자가 "Phase 4 완료, 다음 진행"이라고 승인한 뒤 Phase 5로 이동한다.
3. 실패 시 실패 항목 보고 → 원인 분석 → 수정 → 재검증한다.

---

## 메모 / 주의

- 이 단계에서도 대규모 최적화보다 결정성과 관측성이 우선이다.
- delta snapshot은 복잡해지기 쉽다. full snapshot fallback을 반드시 남긴다.

---
_다음 페이즈: Phase 5 — Blueprint & Reality Bridge Prep → [`wire-craft-phase-5-blueprint-ux-polish.md`](./wire-craft-phase-5-blueprint-ux-polish.md)_
