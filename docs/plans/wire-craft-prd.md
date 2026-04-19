# WireCraft — PRD (MVP Reset)

_Reset date: 2026-04-19_

## Overview

WireCraft는 거대한 오픈월드 샌드박스가 아니라, 초보 메이커가 작은 가상 작업대에서 물리 구조물과 전기/제어 부품을 조합해 반응하는 장치를 만드는 온라인 공작실이다.

기존 Phase 1-4 구현은 서버 권한 루프, 3D 렌더링, 회로 평가, 액추에이터 스냅샷, 기본 협업 동기화의 기술 기반으로 유지한다. Reset 이후 MVP는 이 기반 위에 "사용자가 실제로 장치 하나를 만들 수 있는가"를 검증하는 제품 슬라이스로 축소한다.

## Product Goal

사용자는 작은 작업대 안에서 버튼, 전원, 전선, 제어 모듈, 피스톤/모터, 물리 블록을 조합해 하나의 동작하는 장치를 만든다. 앱은 각 블록의 역할, 방향, 연결 포트, 장착 슬롯, 현재 신호 상태를 보여주고, 사용자가 잘못 배치하거나 위험한 연결을 만들 때 이해 가능한 피드백을 준다.

장기적으로는 자동차, 엘리베이터, 집 배선, 자동문, 로봇 팔 같은 가상 장치를 만들 수 있어야 한다. 하지만 MVP에서는 이 세계를 "작업대 크기의 첫 장치"로 줄인다.

## Primary Persona

어릴 때 라디오 만들기, 과학상자, 만능회로키트, 전구 켜기, 모터 돌리기 같은 것을 해보고 싶었지만 돈, 공간, 장비, 지식 장벽 때문에 못 했고, 이제 온라인에서 안전하게 실험하며 실제 제작 감각을 얻고 싶은 성인 초보 메이커.

## Value Proposition

WireCraft는 회로 시뮬레이터와 복셀 게임 사이의 빈 공간을 겨냥한다. 사용자는 단순한 회로 값을 보는 데서 끝나지 않고, 그 신호가 3D 공간의 물리 부품을 움직이는 것을 본다. 동시에 앱은 실제 하드웨어에서 중요한 pin, port, driver, slot, mounting point 같은 개념을 잃지 않는다.

## Reset MVP: Workbench Device Builder

MVP의 현재 이름은 **Workbench Device Builder**다.

첫 번째 목표 장치는 **버튼으로 여닫는 슬라이딩 도어**다. 이 장치는 자동차나 엘리베이터보다 작지만, WireCraft의 핵심을 모두 포함한다.

- 물리 구조: frame, wall, sliding door panel
- 전기/신호: power, button, wire
- 제어: control board 또는 driver module
- 액추에이터: piston 또는 linear actuator
- 피드백: door open/closed, signal high/low, blocked 상태

## MVP Functional Requirements

| # | 기능 | 설명 | 수용 기준 |
|---|---|---|---|
| 1 | 작은 작업대 월드 | 기본 월드는 거대한 공간이 아니라 제한된 workbench 영역이다. | 사용자는 한 화면에서 첫 장치 전체를 볼 수 있다. |
| 2 | Block instance model | 배치된 블록은 type뿐 아니라 facing, properties, ports, slots, attachments를 가진다. | 배치된 블록을 선택하면 인스턴스 상태를 inspect panel에서 볼 수 있다. |
| 3 | Selection/property editing | 사용자는 배치된 블록을 선택, 회전, 속성 편집할 수 있다. | 방향이 중요한 블록의 facing 변경이 서버 상태와 렌더링에 반영된다. |
| 4 | Port-aware connections | 회로 연결은 단순 인접이 아니라 방향과 port 호환성으로 결정된다. | wire/board/actuator가 올바른 port끼리 연결될 때만 신호가 전달된다. |
| 5 | Module/board attachment | 물리 블록은 제한된 slot에 board/module을 장착할 수 있다. | host block에 control board를 장착하면 host의 동작 규칙이 확장된다. |
| 6 | Sliding door vertical slice | 버튼 신호가 control/driver를 거쳐 door actuator를 움직인다. | 버튼을 누르면 문이 열리고, 떼면 닫히며, 막힌 경우 blocked 상태가 보인다. |
| 7 | Guided first mission | blank canvas 대신 첫 장치 만들기 흐름을 제공한다. | 사용자는 순서대로 배치/연결/테스트하며 첫 장치를 완성할 수 있다. |

## Domain Model Direction

MVP는 다음 개념을 명확히 분리한다.

- **PhysicalBlock**: 공간을 차지하는 구조물. 예: wall, frame, door panel, housing.
- **ComponentBlock**: 신호나 전기 역할을 갖는 부품. 예: power, button, wire, lamp, driver.
- **ModuleBoard**: host block의 slot에 장착되는 제어/회로 모듈.
- **Port**: 방향과 signal type이 있는 연결점. 예: `IN`, `OUT`, `POWER`, `GND`, `ACTUATOR`.
- **Slot**: 특정 module이나 board를 장착할 수 있는 자리. 예: `control_board_slot`.
- **Device**: 여러 block instance, attachments, internal connections를 묶은 작은 장치 단위.

사용자가 말한 "계산기 블록 + 계산기 로직 기판 = 계산 로직이 애드온된 계산기 블록"은 장기적으로 이 모델의 자연스러운 확장이다. MVP에서는 이 기능을 일반화하지 않고, door/lamp 같은 단일 host block과 control board 한 개 수준으로 제한한다.

## Non-Goals for Reset MVP

- 무한 월드, 지형, 채굴, 생존 게임 루프
- 자동차 drivetrain, 완전한 엘리베이터 층 제어, 복잡한 로봇 관절
- 임의 중첩 가능한 블록 그룹/프리팹 시스템
- 완전한 rigid body physics, friction, torque, gear simulation
- KiCad 수준의 회로 설계/검증
- 완전한 Arduino/C++ 런타임
- public multiplayer, account, lobby, marketplace
- 완성형 3D 프린팅/Reality Pack export

## Existing Foundation Kept

- Server-authoritative command/snapshot loop
- WebSocket room/presence foundation
- Snapshot interpolation
- Basic voxel placement/removal
- Starter digital circuit evaluation
- Button input command
- Actuator transform snapshots
- Basic component cards
- Camera navigation and placement facing

## Success Criteria

- 사용자가 작은 작업대에서 첫 sliding door 장치를 완성할 수 있다.
- 사용자는 배치된 블록을 선택하고, 방향과 속성을 편집할 수 있다.
- wire/board/actuator 연결은 port/slot 규칙을 따른다.
- 버튼 입력이 서버 권한 시뮬레이션에서 door motion으로 이어진다.
- inspect panel이 선택한 인스턴스의 role, ports, slots, signal, warnings를 설명한다.
- 모든 핵심 상태는 나중에 blueprint, BOM, wiring guide로 확장할 수 있는 metadata를 보존한다.

## Primary Risks

1. 다시 범위를 자동차/엘리베이터/오픈월드로 넓히면 첫 사용 가능한 장치가 늦어진다.
2. port/slot 모델을 너무 일반화하면 MVP 전에 추상화만 커진다.
3. 단순 인접 회로 모델을 오래 유지하면 방향성 전선과 기판 장착 UX가 막힌다.
4. UI가 배치된 인스턴스가 아니라 부품 카탈로그만 설명하면 사용자는 장치를 고칠 수 없다.

## Planning Links

- Reset plan: [`wire-craft-mvp-reset.md`](./wire-craft-mvp-reset.md)
- Roadmap: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)
- Research notes: [`wire-craft-research-notes.md`](./wire-craft-research-notes.md)
