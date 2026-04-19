# WireCraft — MVP Reset Plan

_Created: 2026-04-19_
_Tracking issue: [#49](https://github.com/devlikebear/wirecraft/issues/49)_

## Why Reset

WireCraft의 기존 구현은 서버 권한 루프, 회로 평가, 액추에이터 snapshot, 멀티플레이 동기화의 기술 기반을 빠르게 검증했다. 하지만 사용자 관점에서는 아직 "장치 하나를 만들고 고치는 경험"이 부족하다.

현재 모델은 대부분 `position + block type + facing`에 머문다. 회로는 모든 인접 블록을 하나의 `body` 핀으로 연결하고, UI는 배치된 인스턴스가 아니라 선택된 부품 카드 중심이다. 이 상태에서 generic collision, metrics, blueprint export를 계속 붙이면 사용 가능한 제작 경험보다 엔진 기능이 먼저 커진다.

Reset의 목적은 기능을 버리는 것이 아니라, 모든 다음 작업을 **작은 작업대에서 첫 동작 장치를 만드는 것**에 맞춰 재정렬하는 것이다.

## Product North Star

사용자가 가상의 작은 작업대에서 물리 구조물, 전기 배선, 제어 모듈, 액추에이터를 조합해 반응하는 장치를 만들고, 왜 움직이는지 inspect panel로 이해할 수 있다.

## Immediate MVP

**Workbench Device Builder**

첫 vertical slice는 **button-controlled sliding door**다.

### Required User Loop

1. 작업대에서 door frame과 door panel을 배치한다.
2. button, power, wire, control/driver module을 배치한다.
3. actuator를 door panel에 연결한다.
4. 배치된 블록을 선택해 facing, properties, ports, slots를 확인한다.
5. 필요한 블록을 회전하거나 속성을 수정한다.
6. 버튼을 눌러 door open/close 동작을 확인한다.
7. 막힌 경우 blocked 상태와 이유를 확인한다.

## Scope Reduction

### Keep

- Server-authoritative command/snapshot model
- WebSocket room foundation
- Snapshot interpolation
- Existing block placement/removal path
- Digital High/Low starter circuit runtime
- Button input
- Kinematic actuator snapshots
- Component card content
- Camera navigation and placement facing

### Pause

- Generic Phase 4 collision work not tied to a concrete device
- Generic metrics work
- Blueprint/export work before block instance metadata is useful
- 4-client multiplayer checkpoint

### Remove From Immediate MVP

- Open-world scale building
- Arbitrary nested block groups
- Full vehicle or elevator systems
- Full PCB/EDA modeling
- Full rigid body physics
- Public rooms/accounts

## Conceptual Model

### Block Instance

A placed block is not only a type at a coordinate. It is an instance with state.

- `id` or deterministic instance key
- `blockType`
- `position`
- `facing`
- `properties`
- `ports`
- `slots`
- `attachments`

The implementation may keep position as the first instance key for compatibility, but the model should stop treating block type as the whole object.

### Port

A port is a directional connection point with a signal kind.

Examples:

- `OUT digital`
- `IN digital`
- `POWER`
- `GND`
- `ACTUATOR`

The circuit graph should eventually connect compatible ports, not all adjacent block bodies.

### Slot

A slot is a constrained attachment point.

Examples:

- `control_board_slot`
- `driver_slot`
- `lamp_module_slot`

MVP allows a simple one-level attachment. A host block may hold one supported module. Arbitrary nested assemblies are out of scope.

### Device

A device is a small named group of block instances, attachments, and connections. Device save/load comes after the selected block and port/slot model exists.

## First Device: Sliding Door

### Parts

- `solid` or `frame`
- `door_panel`
- `power`
- `button`
- `wire`
- `control_board` or `motor_driver`
- `piston` or `linear_actuator`

### Behavior

- Released button: signal low, door closed.
- Pressed button: signal high, actuator extends, door opens.
- Blocked path: actuator stops or clamps, `blocked` property becomes visible.

### Why This Device

- It is smaller than a car or elevator.
- It still requires physical structure, wiring, control, actuator motion, and feedback.
- It makes block direction, ports, slots, and collision meaningful immediately.

## Revised Next Work

The next implementation should not be generic collision. It should start with block instance state because all following UX depends on it.

Recommended next work order:

**WO-R1: Add block instance properties to snapshots** ([#50](https://github.com/devlikebear/wirecraft/issues/50))

- Add a minimal `properties` map to placed block snapshots.
- Preserve existing placement/removal behavior.
- Add tests that default snapshots include empty properties.
- Do not add property editing UI yet.

This gives the product a place to store later state such as `axis`, `blocked`, `open`, `moduleSlot`, or user-facing labels without redesigning the protocol every time.

## GitHub Issue Reclassification

- [#39](https://github.com/devlikebear/wirecraft/issues/39) remains useful historical context but is paused by this reset.
- [#47](https://github.com/devlikebear/wirecraft/issues/47) is no longer the immediate next step; collision should return inside the sliding door vertical slice.
- [#49](https://github.com/devlikebear/wirecraft/issues/49) is the reset parent issue.

## Completion Criteria For Reset

- PRD names Workbench Device Builder as current MVP.
- Roadmap has Reset phases R0-R5.
- Current status points new sessions to reset work orders.
- Phase 4 docs no longer list #47 as the next implementation task.
- GitHub issues record the reset decision.
