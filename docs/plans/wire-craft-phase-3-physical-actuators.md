# Phase 3: Physical Actuators — 작업지시서

_작성일: 2026-04-18_
_속한 로드맵: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)_
_예상 소요: 1-2주_

## 페이즈 목표

회로 신호가 3D 월드의 물리 블록 움직임을 만들도록 연결한다. MVP에서는 완전한 강체 물리보다 결정적인 kinematic actuator를 우선한다. 이 단계가 끝나면 버튼 또는 MCU output 신호가 피스톤이나 모터 transform을 바꾸고, 클라이언트는 서버 snapshot을 보간해 움직임을 부드럽게 렌더링한다.

## 전제 조건

- [x] Phase 2 완료 및 사용자 승인
- [ ] snapshot interpolation이 Phase 1에서 동작 중
- [ ] circuit output state를 simulation layer에서 읽을 수 있음

## 포함 기능

1. actuator component model
2. piston block
3. motor block 또는 rotating debug block
4. signal-to-motion mapping
5. transform snapshot replication
6. 최소 sensor input stub
7. motor driver/transistor 현실 제약 warning

## 이 페이즈에서 하지 않는 것

- 복잡한 강체 충돌/조인트 → Phase 4 이후
- PWM의 정확한 전기적 해석 → Later
- 다축 로봇/공장 라인 템플릿 → Later

## 작업 체크리스트

### 작업 그룹 A: Actuator model

- [x] **T3.A.1** — dynamic entity와 transform type 정의.
  - 파일:
    - `internal/physics/transform.go`
    - `internal/physics/entity.go`
    - `internal/physics/entity_test.go`
  - 내용:
    - `Vec3`, `Quat`, `Transform`
    - `EntityID`
    - `DynamicEntity{ID, Type, Transform, Velocity, Target}`
    - deterministic update를 위한 entity ordering
  - 검증: `go test ./internal/physics/...`

- [x] **T3.A.2** — actuator component 구현.
  - 파일:
    - `internal/actuator/types.go`
    - `internal/actuator/piston.go`
    - `internal/actuator/piston_test.go`
  - 내용:
    - `ActuatorType`: `piston`, `motor`
    - 입력 signal: Low/High/PWM placeholder
    - piston: High일 때 extension 1 block, Low일 때 retract
    - motor: `motor_driver` 또는 `transistor_switch` 입력을 통해서만 enable
    - movement speed 제한
  - 테스트:
    - High 입력 시 target extension 계산
    - Low 입력 시 target retract 계산
    - tick delta에 따른 위치 보간
  - 검증: `go test ./internal/actuator/...`

### 작업 그룹 B: Circuit-to-actuator integration

- [x] **T3.B.1** — actuator block type 추가.
  - 파일:
    - `internal/world/block.go`
    - `internal/circuit/types.go`
    - `web/src/net/protocol.ts`
  - 내용:
    - block type: `piston`, `motor`, `motor_driver`, `transistor_switch`
    - actuator input pin 연결 규칙 정의
    - GPIO pin에 motor를 직접 연결하면 invalid wiring 또는 warning으로 처리
  - 검증: `go test ./... && cd web && npm run build`

- [x] **T3.B.2** — simulation update 순서 확정.
  - 파일:
    - `internal/sim/simulation.go`
    - `internal/sim/simulation_test.go`
  - 내용:
    - command drain
    - world update
    - circuit evaluation
    - actuator input mapping
    - actuator/physics update
    - snapshot build
  - 테스트:
    - 버튼 High 후 같은 tick 또는 다음 tick에서 actuator target이 갱신되는지 명확히 검증
  - 검증: `go test ./internal/sim/...`

- [x] **T3.B.3** — transform snapshot schema 추가.
  - 파일:
    - `internal/netproto/snapshot.go`
    - `web/src/net/protocol.ts`
  - 내용:
    - `entities`에 actuator transform 포함
    - entity id, type, position, rotation, scale
  - 검증: `go test ./... && cd web && npm run build`

### 작업 그룹 C: Client rendering

- [x] **T3.C.1** — actuator mesh 렌더링.
  - 파일:
    - `web/src/render/EntityRenderer.ts`
    - `web/src/render/ActuatorMeshes.ts`
  - 내용:
    - piston base와 moving head 표현
    - motor는 회전 큐브 또는 축 표시
    - server entity transform을 interpolation 결과로 적용
  - 검증: 브라우저에서 piston/motor 표시 확인

- [x] **T3.C.2** — actuator block placement UI.
  - 파일:
    - `web/src/ui/Toolbar.ts`
    - `web/src/input/EditController.ts`
  - 내용:
    - piston/motor/motor driver/transistor switch 선택
    - orientation은 기존 placement semantics 유지
  - 검증: actuator block 배치 후 서버 snapshot에 entity 생성 확인

### 작업 그룹 D: Sensor input stub

- [x] **T3.D.1** — button 외 sensor input 확장 포인트 작성.
  - 파일:
    - `internal/sensor/types.go`
    - `internal/sensor/types_test.go`
  - 내용:
    - `SensorType`: `button`, `proximity_stub`
    - proximity는 Phase 3에서 단순 거리 조건만 지원하거나 stub로 남김
  - 검증: `go test ./internal/sensor/...`

### 작업 그룹 E: 초보자 현실 제약 표시

- [x] **T3.E.1** — motor/driver component card를 추가한다.
  - 파일:
    - `internal/component/card.go`
    - `docs/reference/component-cards.md`
    - `web/src/state/componentCards.ts`
  - 내용:
    - DC motor는 GPIO에서 직접 구동하지 않는다는 warning
    - motor driver/transistor/diode의 역할 설명
    - Wire Craft 시뮬레이션은 전류/역전압을 단순화한다고 명시
  - 검증: motor 관련 block 선택 시 inspect panel에 warning 표시

---

## ✅ Phase 3 Checkpoint

**구현 확인:**
- [x] 피스톤 또는 모터 block을 배치할 수 있다.
- [ ] 회로 output High/Low가 actuator input으로 연결된다.
- [ ] 서버 tick에서 actuator transform이 계산된다.
- [ ] 클라이언트가 actuator transform을 보간 렌더링한다.
- [x] sensor 확장 지점이 코드상 분리되어 있다.
- [x] GPIO pin에 motor를 직접 연결하면 warning 또는 invalid wiring이 표시된다.

**자동 검증:**
- [ ] 서버 테스트 통과: `go test ./...`
- [ ] 클라이언트 테스트 통과: `cd web && npm test`
- [ ] 클라이언트 빌드 통과: `cd web && npm run build`

**수동 확인:**
- [ ] 버튼 -> 전선 -> 피스톤 회로를 만들고 버튼을 누르면 피스톤이 움직인다.
- [ ] 버튼을 떼면 피스톤이 원위치로 돌아온다.
- [ ] motor block은 driver/transistor component 없이 직접 연결하면 동작하지 않거나 warning을 표시한다.
- [ ] 브라우저 창 2개에서 피스톤 위치가 동일하게 보인다.
- [ ] 움직임이 서버 tick 단위로 끊겨 보이지 않고 보간된다.

**완료 처리:**
1. 위 항목 모두 통과 시 완료 요약과 검증 결과를 보고한다.
2. 사용자가 "Phase 3 완료, 다음 진행"이라고 승인한 뒤 Phase 4로 이동한다.
3. 실패 시 실패 항목 보고 → 원인 분석 → 수정 → 재검증한다.

---

## 메모 / 주의

- Phase 3의 물리는 “정확한 강체 시뮬레이션”보다 “서버 결정성과 네트워크 시각화”가 우선이다.
- actuator 움직임이 흔들리면 interpolation delay, server time estimation, entity id 안정성을 먼저 점검한다.

---
_다음 페이즈: Phase 4 — Multiplayer Physics Sync → [`wire-craft-phase-4-multiplayer-physics-sync.md`](./wire-craft-phase-4-multiplayer-physics-sync.md)_
