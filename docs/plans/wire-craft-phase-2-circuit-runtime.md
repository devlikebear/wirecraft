# Phase 2: Circuit Runtime — 작업지시서

_작성일: 2026-04-18_
_속한 로드맵: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)_
_예상 소요: 1-2주_

## 페이즈 목표

서버가 복셀 월드 안의 전자 블록을 회로 그래프로 해석하고, tick마다 High/Low 신호를 결정적으로 계산한다. 이 단계가 끝나면 사용자는 버튼, 전원, 전선, AND 게이트, MCU output node를 배치해 간단한 논리 회로를 만들고 모든 클라이언트에서 같은 결과를 볼 수 있다.

## 전제 조건

- [x] Phase 1 완료 및 사용자 승인
- [x] snapshot schema에 block state 또는 circuit debug state를 추가할 수 있는 구조 확보

## 포함 기능

1. 회로 블록 타입 추가
2. world block에서 circuit graph 추출
3. signal propagation
4. 버튼 입력 command
5. 회로 상태 시각화
6. 초보자용 component card

## 이 페이즈에서 하지 않는 것

- 액추에이터 움직임 → Phase 3
- 복잡한 아날로그 전기 시뮬레이션 → Out of Scope
- Arduino/C++ 호환 런타임 → Later

## 작업 체크리스트

### 작업 그룹 A: Circuit domain model

- [ ] **T2.A.1** — 회로 블록 타입 정의.
  - 파일:
    - `internal/world/block.go`
    - `internal/circuit/types.go`
  - 내용:
    - block type: `power`, `wire`, `button`, `and_gate`, `mcu_output`
    - 방향이 필요한 블록의 orientation 필드 추가
    - block metadata 구조 추가
  - 검증: `go test ./internal/world/...`

- [ ] **T2.A.2** — circuit graph 자료구조 구현.
  - 파일:
    - `internal/circuit/graph.go`
    - `internal/circuit/graph_test.go`
  - 내용:
    - `NodeID`, `NodeType`, `SignalState`
    - graph node와 edge
    - 인접한 wire/pin 연결 규칙
  - 테스트:
    - 전원-전선-출력 연결 생성
    - 방향성 게이트 pin 연결 생성
  - 검증: `go test ./internal/circuit/...`

- [ ] **T2.A.3** — component card 데이터 모델을 만든다.
  - 파일:
    - `internal/component/card.go`
    - `internal/component/card_test.go`
    - `docs/reference/component-cards.md`
  - 내용:
    - 필드: `id`, `displayName`, `whatItDoes`, `howToWire`, `realWorldWarning`, `simulationSimplification`, `pins`
    - 초기 카드: LED, resistor, power, ground, wire, button, PWM pin, MCU output
    - `docs/reference/component-cards.md`에 초보자 설명 초안 작성
  - 테스트:
    - 필수 필드 누락 검증
    - pin capability 정의 검증
  - 검증: `go test ./internal/component/...`

### 작업 그룹 B: 회로 평가

- [ ] **T2.B.1** — world에서 circuit graph를 추출한다.
  - 파일:
    - `internal/circuit/extract.go`
    - `internal/circuit/extract_test.go`
  - 내용:
    - world snapshot을 순회해 회로 관련 block만 추출
    - 인접 규칙으로 edge 생성
    - dirty block 변경 시 재추출할 수 있는 구조로 작성
  - 검증: `go test ./internal/circuit/...`

- [ ] **T2.B.2** — digital signal evaluation 구현.
  - 파일:
    - `internal/circuit/evaluate.go`
    - `internal/circuit/evaluate_test.go`
  - 내용:
    - High/Low/Unknown 상태
    - power source propagation
    - button state 반영
    - AND gate truth table
    - evaluation order가 deterministic하도록 sort 적용
  - 테스트:
    - power -> wire -> output High
    - button off -> output Low
    - AND gate 4가지 truth table
    - cycle이 있어도 panic 없이 안정 상태 또는 Unknown 처리
  - 검증: `go test ./internal/circuit/...`

### 작업 그룹 C: 서버 tick 통합

- [ ] **T2.C.1** — simulation tick에 circuit evaluation을 통합한다.
  - 파일:
    - `internal/sim/simulation.go`
    - `internal/sim/simulation_test.go`
  - 내용:
    - command drain
    - world update
    - circuit graph extraction/evaluation
    - snapshot build 순서 고정
  - 검증: `go test ./internal/sim/...`

- [ ] **T2.C.2** — 버튼 입력 command 구현.
  - 파일:
    - `internal/netproto/command.go`
    - `internal/sim/commands.go`
  - 내용:
    - command type: `set_button`
    - button press/release 상태 저장
    - clientId와 commandId 검증 유지
  - 검증: `go test ./...`

- [ ] **T2.C.3** — snapshot에 circuit state 추가.
  - 파일:
    - `internal/netproto/snapshot.go`
    - `web/src/net/protocol.ts`
  - 내용:
    - block position별 signal state 또는 circuit debug state 포함
    - Phase 2는 이해하기 쉬운 구조 우선, 압축은 Phase 4로 이동
  - 검증:
    - `go test ./...`
    - `cd web && npm run build`

### 작업 그룹 D: 클라이언트 회로 표시

- [ ] **T2.D.1** — 회로 블록 배치 UI를 추가한다.
  - 파일:
    - `web/src/ui/Toolbar.ts`
    - `web/src/input/EditController.ts`
  - 내용:
    - solid, power, wire, button, and_gate, mcu_output 선택
    - 선택된 block type을 place command에 포함
  - 검증: 브라우저에서 각 블록 배치 확인

- [ ] **T2.D.2** — signal state visualization 구현.
  - 파일:
    - `web/src/render/VoxelRenderer.ts`
    - `web/src/render/CircuitOverlay.ts`
  - 내용:
    - High 상태는 밝은 emissive material 또는 overlay
    - Low 상태는 muted material
    - Unknown은 debug color
  - 검증: 버튼 상태 변경에 따라 overlay 갱신 확인

- [ ] **T2.D.3** — inspect panel에서 component card를 표시한다.
  - 파일:
    - `web/src/ui/InspectPanel.ts`
    - `web/src/state/componentCards.ts`
  - 내용:
    - 선택한 부품의 역할, 배선법, 현실 주의사항, 시뮬레이션 단순화 범위 표시
    - High/Low/PWM 같은 용어는 tooltip으로 설명
  - 검증: 브라우저에서 부품 선택 후 card 내용 표시 확인

---

## ✅ Phase 2 Checkpoint

**구현 확인:**
- [ ] 전원, 전선, 버튼, AND 게이트, MCU output block type이 배치된다.
- [ ] 서버 tick마다 회로 상태가 결정적으로 평가된다.
- [ ] 버튼 입력이 snapshot의 signal state에 반영된다.
- [ ] 클라이언트가 High/Low 상태를 시각적으로 표시한다.
- [ ] 선택한 부품의 component card가 inspect panel에 표시된다.

**자동 검증:**
- [ ] 서버 테스트 통과: `go test ./...`
- [ ] 클라이언트 빌드 통과: `cd web && npm run build`
- [ ] 회로 truth table 테스트 통과: `go test ./internal/circuit/...`

**수동 확인:**
- [ ] power -> wire -> mcu_output 회로를 만들면 output이 High로 표시된다.
- [ ] button + power -> wire -> mcu_output 회로에서 버튼을 누를 때만 High가 표시된다.
- [ ] AND gate 두 입력이 모두 High일 때만 출력이 High로 표시된다.
- [ ] 브라우저 창 2개에서 회로 상태가 동일하게 보인다.
- [ ] 초보자가 LED/button/PWM pin card를 보고 각 부품의 역할과 주의사항을 확인할 수 있다.

**완료 처리:**
1. 위 항목 모두 통과 시 완료 요약과 검증 결과를 보고한다.
2. 사용자가 "Phase 2 완료, 다음 진행"이라고 승인한 뒤 Phase 3로 이동한다.
3. 실패 시 실패 항목 보고 → 원인 분석 → 수정 → 재검증한다.

---

## 메모 / 주의

- 이 단계는 디지털 논리만 다룬다. 저항, 전류량, 전압강하 같은 실제 전기 해석은 MVP 범위가 아니다.
- circuit evaluation은 테스트 우선으로 작성한다. 나중에 물리와 결합되면 디버깅이 어려워진다.

---
_다음 페이즈: Phase 3 — Physical Actuators → [`wire-craft-phase-3-physical-actuators.md`](./wire-craft-phase-3-physical-actuators.md)_
