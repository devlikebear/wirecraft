# WireCraft — 개발 로드맵 (MVP Reset)

_Reset date: 2026-04-19_
_Tracking issue: [#49](https://github.com/devlikebear/wirecraft/issues/49)_

## Overview

WireCraft의 장기 비전은 가상 세계에서 물리 구조물, 전기 배선, 제어 로직, 액추에이터를 조합해 실제처럼 반응하는 장치를 만드는 온라인 공작실이다. 다만 첫 MVP는 마인크래프트 규모의 월드가 아니다. Reset 이후 MVP는 **Workbench Device Builder**로 축소한다.

현재 목표는 작은 작업대에서 **버튼으로 여닫는 슬라이딩 도어**를 만드는 것이다. 이 장치는 물리 블록, 전선, 버튼, 제어/드라이버 모듈, 액추에이터, 충돌/blocked 상태, inspect UI를 모두 포함하는 가장 작은 제품 슬라이스다.

## Reset Decision

기존 Phase 1-3과 Phase 4 일부 작업은 버리지 않는다. 서버 권한 루프, WebSocket room, command ordering, changed-set snapshot, 회로 평가, 액추에이터 snapshot, camera control은 유지한다.

다만 Phase 4의 남은 generic collision/metrics 작업과 Phase 5의 blueprint/export 작업은 즉시 진행하지 않는다. 이 작업들은 첫 장치 모델과 연결된 뒤 다시 진행한다.

## Current MVP Completion Criteria

- [ ] 사용자가 제한된 workbench 영역에서 첫 장치 전체를 볼 수 있다.
- [ ] 배치된 block instance가 type, position, facing, properties를 가진다.
- [ ] 배치된 블록을 선택해 inspect panel에서 인스턴스 상태를 볼 수 있다.
- [ ] 방향과 port/slot 규칙에 따라 wire, board, actuator 연결이 결정된다.
- [ ] 물리 block에 module/board를 제한적으로 장착할 수 있다.
- [ ] 버튼 신호가 control/driver를 거쳐 sliding door actuator를 움직인다.
- [ ] actuator가 막히면 blocked 상태가 표시된다.
- [ ] guided mission으로 첫 sliding door 장치를 완성할 수 있다.
- [ ] 이후 blueprint/BOM/wiring guide로 확장할 metadata를 잃지 않는다.

## Phase History

### Foundation 1: Authoritative Voxel Loop

Status: Completed. Tracking issue: [#1](https://github.com/devlikebear/wirecraft/issues/1).

서버 권한 tick loop, WebSocket command/snapshot protocol, voxel rendering, raycast edit, interpolation foundation을 확보했다.

### Foundation 2: Circuit Runtime

Status: Completed. Tracking issue: [#17](https://github.com/devlikebear/wirecraft/issues/17).

전원, 전선, 버튼, AND gate, MCU output의 최소 digital signal evaluation과 component card foundation을 확보했다.

### Foundation 3: Physical Actuators

Status: Completed. Tracking issue: [#28](https://github.com/devlikebear/wirecraft/issues/28).

회로 신호가 kinematic piston/motor entity snapshot으로 이어지는 초기 actuator foundation을 확보했다.

### Foundation 4: Multiplayer Physics Sync

Status: Paused by MVP Reset. Tracking issue: [#39](https://github.com/devlikebear/wirecraft/issues/39).

room/session model, presence, deterministic command ordering, command acknowledgement, changed-set snapshot, viewport navigation, placement facing까지 완료했다. 남은 generic collision/metrics/checkpoint는 첫 device slice에 맞게 재분류한다.

## Reset Phases

### Phase R0: MVP Reset Planning

Status: Completed locally. Tracking issue: [#49](https://github.com/devlikebear/wirecraft/issues/49).

- **목표**: 현재 구현과 기획을 Workbench Device Builder 기준으로 재정렬한다.
- **산출물**: reset PRD, reset roadmap, reset work orders, current status update.
- **완료 조건**: 다음 구현 작업이 작고 테스트 가능한 work order로 정리되어 있다.

### Phase R1: Workbench Block Instance UX

- **목표**: "블록 타입을 놓는 앱"에서 "배치된 블록 인스턴스를 선택하고 편집하는 앱"으로 전환한다.
- **포함 기능**: block properties, selected block state, rotate existing block, inspect instance panel, workbench bounds.
- **완료된 첫 작업**: [`WO-R1: Add block instance properties to snapshots`](https://github.com/devlikebear/wirecraft/issues/50).
- **다음 작업**: [`WO-R2: Select placed block and show instance state`](https://github.com/devlikebear/wirecraft/issues/51).
- **Checkpoint**: 배치된 wire/button/piston을 선택하면 position, facing, editable properties, component metadata가 inspect panel에 표시된다.

### Phase R2: Port and Attachment Model

- **목표**: 단순 인접 회로를 방향/포트/슬롯 기반 연결 모델로 전환한다.
- **포함 기능**: port definitions, compatible signal types, slot definitions, one-board attachment rule, invalid/warning feedback.
- **Checkpoint**: wire와 board가 호환되는 port끼리 연결될 때만 신호가 전달된다.

### Phase R3: Sliding Door Vertical Slice

- **목표**: 첫 실제 장치인 button-controlled sliding door를 완성한다.
- **포함 기능**: door frame/panel, control board/driver, piston axis from facing, blocked state, door open/closed state.
- **Checkpoint**: 버튼을 누르면 문이 열리고, 떼면 닫히며, 막힌 경우 blocked 상태가 inspect panel에 표시된다.

### Phase R4: Guided Mission and Device Save

- **목표**: blank canvas 대신 첫 장치를 따라 만들 수 있는 흐름을 제공하고, 작은 device 단위 저장을 준비한다.
- **포함 기능**: guided mission steps, build log, device JSON schema, local save/load.
- **Checkpoint**: 사용자가 안내를 따라 sliding door를 만들고, device JSON으로 저장/불러올 수 있다.

### Phase R5: Multiplayer and Reality Bridge Re-entry

- **목표**: 제품 루프가 쓸만해진 뒤 기존 Phase 4/5의 협업/blueprint/BOM/wiring/export 작업을 다시 연결한다.
- **포함 기능**: 2-client device collaboration smoke, server metrics, BOM prototype, wiring guide prototype.
- **Checkpoint**: 두 클라이언트가 같은 device를 보고, 저장된 device에서 wiring metadata를 추출할 수 있다.

## Dependency Order

```text
Completed foundation:
  authoritative loop -> circuit runtime -> actuator snapshots -> partial multiplayer sync

Reset path:
  R0 planning
    -> R1 block instance UX
      -> R2 port/slot/attachment model
        -> R3 sliding door device
          -> R4 guided mission/device save
            -> R5 multiplayer/reality bridge re-entry
```

## Out of Scope Until R5+

- 자동차 drivetrain
- 엘리베이터 전체 층 제어
- 임의 중첩 가능한 block group system
- full rigid body physics
- public multiplayer
- Reality Pack export package
- Arduino-compatible runtime

## Active Work Order Source

- Reset work orders: [`../tasks/mvp-reset-work-orders.md`](../tasks/mvp-reset-work-orders.md)
- Current status: [`../tasks/current-status.md`](../tasks/current-status.md)
- PRD: [`wire-craft-prd.md`](./wire-craft-prd.md)
- Reset rationale: [`wire-craft-mvp-reset.md`](./wire-craft-mvp-reset.md)
