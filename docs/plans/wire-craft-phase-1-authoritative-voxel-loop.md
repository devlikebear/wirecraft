# Phase 1: Authoritative Voxel Loop — 작업지시서

_작성일: 2026-04-18_
_속한 로드맵: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)_
_예상 소요: 1-2주_

## 페이즈 목표

서버 권한 구조의 가장 작은 end-to-end 루프를 만든다. Go 서버가 world state를 고정 tick으로 관리하고, Three.js 클라이언트는 snapshot을 받아 복셀 월드를 렌더링하며, 사용자의 블록 배치/삭제는 command로만 서버에 전송된다. 이 페이즈에서 snapshot interpolation까지 넣어 서버 20Hz와 클라이언트 60FPS 간극을 초기에 검증한다.

## 전제 조건

- [ ] 새 Git repository 초기화 여부 결정
- [ ] Go 1.22+ 설치
- [ ] Node.js 20+ 설치
- [ ] package manager 선택: npm 또는 pnpm

## 포함 기능

1. 프로젝트 기본 구조
2. Go WebSocket 서버와 fixed tick loop
3. voxel world 자료구조와 command validation
4. snapshot protocol
5. Vite + Three.js 클라이언트
6. raycast 기반 block place/remove
7. snapshot buffer interpolation
8. debug overlay

## 이 페이즈에서 하지 않는 것

- 회로 엔진 → Phase 2
- 액추에이터/센서 → Phase 3
- delta compression → Phase 4
- blueprint 저장 → Phase 5
- 계정/권한/DB → Out of Scope

## 작업 체크리스트

### 작업 그룹 A: 프로젝트 초기화

- [ ] **T1.A.1** — repository와 기본 디렉터리 구조를 만든다.
  - 파일:
    - `go.mod`
    - `cmd/wirecraft-server/main.go`
    - `internal/world/`
    - `internal/netproto/`
    - `web/`
    - `README.md`
  - 내용:
    - Go module name 결정
    - 서버와 클라이언트 분리 구조 확정
    - README에 MVP 범위와 실행 명령 초안 작성
  - 검증: `go test ./...`

- [ ] **T1.A.2** — Vite + TypeScript + Three.js 클라이언트 초기화.
  - 파일:
    - `web/package.json`
    - `web/src/main.ts`
    - `web/src/App.ts`
    - `web/src/styles.css`
  - 내용:
    - Three.js scene, camera, renderer, resize handling 구성
    - 빈 grid helper와 basic lighting 표시
  - 검증:
    - `cd web && npm install`
    - `cd web && npm run build`

### 작업 그룹 B: 서버 tick loop와 world model

- [ ] **T1.B.1** — fixed tick game loop 구현.
  - 파일:
    - `internal/sim/tick.go`
    - `internal/sim/tick_test.go`
  - 내용:
    - `type TickID uint64`
    - `type Clock struct { RateHz int }`
    - command queue drain -> simulation update -> snapshot publish 순서 정의
    - tick duration 측정 hook 추가
  - 테스트:
    - 20Hz 설정 시 tick duration target이 50ms로 계산되는지
    - tick id가 단조 증가하는지
  - 검증: `go test ./internal/sim/...`

- [ ] **T1.B.2** — voxel world 자료구조 구현.
  - 파일:
    - `internal/world/world.go`
    - `internal/world/world_test.go`
  - 내용:
    - `Position{X,Y,Z int}`
    - `BlockType` enum: `air`, `solid`, `debug_mover`
    - `World.Get`, `World.Set`, `World.Remove`
    - 월드 크기 제한: 기본 32x32x16
  - 테스트:
    - 범위 밖 좌표 거부
    - 같은 좌표 set/remove 결과 결정성 확인
  - 검증: `go test ./internal/world/...`

- [ ] **T1.B.3** — client command 모델 구현.
  - 파일:
    - `internal/netproto/command.go`
    - `internal/netproto/command_test.go`
  - 내용:
    - command type: `place_block`, `remove_block`
    - 필드: `clientId`, `commandId`, `tickHint`, `position`, `blockType`
    - validation: 좌표 범위, block type, command id 중복
  - 검증: `go test ./internal/netproto/...`

### 작업 그룹 C: WebSocket protocol과 snapshot

- [ ] **T1.C.1** — WebSocket server endpoint 구현.
  - 파일:
    - `internal/server/server.go`
    - `internal/server/ws.go`
    - `cmd/wirecraft-server/main.go`
  - 내용:
    - `GET /ws` endpoint
    - 연결 시 임시 `clientId` 발급
    - command 수신 후 simulation queue에 전달
    - tick snapshot broadcast
  - 검증: `go test ./internal/server/...`

- [ ] **T1.C.2** — snapshot schema 구현.
  - 파일:
    - `internal/netproto/snapshot.go`
    - `internal/netproto/snapshot_test.go`
  - 내용:
    - `tick`, `serverTimeMs`, `blocks`, `entities`, `stats`
    - Phase 1은 full snapshot 허용
    - dynamic debug entity 1개를 매 tick 왕복 이동시켜 interpolation 테스트에 사용
  - 검증: `go test ./internal/netproto/...`

### 작업 그룹 D: Three.js 렌더링과 편집 command

- [ ] **T1.D.1** — WebSocket client와 snapshot store 구현.
  - 파일:
    - `web/src/net/socket.ts`
    - `web/src/net/protocol.ts`
    - `web/src/state/snapshotStore.ts`
  - 내용:
    - 서버 연결/재연결 최소 처리
    - snapshot parse와 buffer append
    - command send helper
  - 검증: `cd web && npm test` 또는 `cd web && npm run build`

- [ ] **T1.D.2** — voxel InstancedMesh renderer 구현.
  - 파일:
    - `web/src/render/VoxelRenderer.ts`
  - 내용:
    - block type별 material
    - full snapshot blocks를 InstancedMesh matrix로 변환
    - 최대 block count 제한과 rebuild path 명시
  - 검증: `cd web && npm run build`

- [ ] **T1.D.3** — raycast 기반 배치/삭제 구현.
  - 파일:
    - `web/src/input/EditController.ts`
    - `web/src/main.ts`
  - 내용:
    - 좌클릭: block place
    - Shift+좌클릭 또는 우클릭: block remove
    - 클라이언트는 optimistic write를 하지 않고 서버 snapshot을 기다림
  - 검증:
    - 브라우저에서 클릭 후 서버 snapshot 반영 확인

### 작업 그룹 E: Interpolation

- [ ] **T1.E.1** — snapshot interpolation 순수 로직 구현.
  - 파일:
    - `web/src/sim/interpolation.ts`
    - `web/src/sim/interpolation.test.ts`
  - 내용:
    - `findSnapshotPair(buffer, renderServerTimeMs)`
    - `interpolateTransform(before, after, alpha)`
    - position은 lerp, rotation은 quaternion slerp
    - `interpolationDelayMs` 기본값 120ms
  - 테스트:
    - 두 snapshot 사이 alpha 계산
    - after snapshot 누락 시 fallback
    - buffer 정리 정책
  - 검증: `cd web && npm test`

- [ ] **T1.E.2** — render loop에 interpolation 적용.
  - 파일:
    - `web/src/render/EntityRenderer.ts`
    - `web/src/main.ts`
  - 내용:
    - estimated server time 계산
    - dynamic debug entity에 interpolation 결과 적용
    - static voxel update와 dynamic entity update 분리
  - 검증:
    - 서버 20Hz, 브라우저 60FPS에서 debug entity 움직임 육안 확인

- [ ] **T1.E.3** — debug overlay 구현.
  - 파일:
    - `web/src/ui/DebugOverlay.ts`
    - `web/src/styles.css`
  - 내용:
    - FPS
    - server tick
    - snapshot buffer length
    - interpolation delay
    - RTT 추정값
    - dropped/late snapshot count
  - 검증: 브라우저에서 overlay 값 갱신 확인

---

## ✅ Phase 1 Checkpoint

**구현 확인:**
- [ ] Go 서버가 20Hz fixed tick으로 실행된다.
- [ ] 브라우저 클라이언트가 WebSocket으로 서버에 연결된다.
- [ ] 복셀 블록 배치/삭제가 서버 snapshot을 통해 반영된다.
- [ ] dynamic debug entity가 snapshot interpolation으로 렌더링된다.
- [ ] debug overlay에서 tick/buffer/FPS 관련 값이 보인다.

**자동 검증:**
- [ ] 서버 테스트 통과: `go test ./...`
- [ ] 클라이언트 테스트 통과: `cd web && npm test`
- [ ] 클라이언트 빌드 통과: `cd web && npm run build`

**수동 확인:**
- [ ] 브라우저 창 2개를 열고 같은 서버에 연결하면 동일한 복셀 배치 결과가 보인다.
- [ ] 한 창에서 블록을 배치하면 다른 창에 다음 snapshot 이후 표시된다.
- [ ] debug mover가 서버 tick보다 부드럽게 움직인다.

**완료 처리:**
1. 위 항목 모두 통과 시 Claude Code는 완료 요약, 테스트 결과, 수동 확인 항목을 보고한다.
2. 사용자가 "Phase 1 완료, 다음 진행"이라고 승인한 뒤 Phase 2로 이동한다.
3. 실패 시 실패 항목 보고 → 원인 분석 → 수정 → 재검증한다.

---

## 메모 / 주의

- interpolation은 “나중에 최적화”가 아니라 이 프로젝트의 성립 조건이다. Phase 1에서 반드시 구현한다.
- 클라이언트는 서버 authoritative state를 직접 수정하지 않는다.
- Phase 1 snapshot은 full snapshot으로 시작한다. delta compression은 Phase 4에서 한다.

---
_다음 페이즈: Phase 2 — Circuit Runtime → [`wire-craft-phase-2-circuit-runtime.md`](./wire-craft-phase-2-circuit-runtime.md)_

