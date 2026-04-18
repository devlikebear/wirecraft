# Wire Craft — PRD (MVP)

_작성일: 2026-04-18_

## Overview

Wire Craft는 서버 권한(Server-Authoritative) 구조를 기반으로, 여러 사용자가 하나의 3D 복셀 작업 공간에서 전자 회로와 물리 장치를 함께 만들고 실험하는 협업형 엔지니어링 샌드박스다. 기존 레드스톤류의 추상 논리 회로와 달리 GPIO, PWM, 센서 입력, 액추에이터 출력 같은 하드웨어 개념을 복셀 월드의 모터, 피스톤, 회전축, 충돌체와 직접 연결하는 것이 핵심 차별점이다.

MVP는 완성형 게임이 아니라 핵심 기술 리스크를 검증하는 playable prototype이다. 첫 버전의 목표는 “서버가 모든 상태를 계산하고, 클라이언트는 스냅샷을 보간 렌더링하며, 사용자는 간단한 회로 신호로 물리 블록 하나를 움직일 수 있다”까지다.

## 제품 동기

이 프로젝트의 출발점은 “어릴 때 라디오 만들기, 과학상자, 만능회로키트, 전구 켜기, 모터 돌리기 같은 것을 해보고 싶었지만 돈과 환경 때문에 못 했던 경험”이다. 어른이 된 뒤에도 실제 로봇 제작은 공간, 부품비, 장비, 안전, 전자/기계 지식 장벽이 높다.

Wire Craft는 이 장벽을 온라인에서 낮추는 것을 목표로 한다. 사용자는 가상 공간에서 마음껏 배선하고, 실패하고, 고치고, 움직이는 장치를 만든다. 이후에는 그 결과를 부품 목록, 배선표, 코드, 3D 출력 후보 파일로 연결해 현실 제작으로 넘어갈 수 있어야 한다.

## 타겟 사용자

- **Primary**: 어릴 때 메이커 키트와 로봇 제작을 해보고 싶었지만 돈, 공간, 장비, 지식 장벽 때문에 못 했고, 이제 온라인에서 안전하게 실험하며 실제 제작으로 넘어가고 싶은 성인 초보 메이커
- **Secondary**: 회로/로봇/코딩을 학생에게 가르치고 싶은 교사
- **Secondary**: 3D 프린터와 아두이노를 갖고 있지만 설계-배선-코드 연결이 막히는 취미 제작자
- **Secondary**: 협업으로 작은 자동화 장치, 공장 라인, 로봇 프로토타입을 만들고 싶은 소규모 팀

## 가치 제안

하드웨어와 물리 장치의 관계를 3D 공간에서 직접 실험하고 싶은 초보 메이커를 위한 온라인 공작실로, 서버 권한 틱 루프와 회로/물리 하이브리드 엔진 덕분에 기존 단일 회로 시뮬레이터나 레드스톤식 추상 회로보다 더 실제적인 엔지니어링 피드백을 제공한다. 장기적으로는 가상 제작물을 실제 부품 목록, 배선표, 컨트롤러 코드, 3D 출력 파일로 이어주는 Reality Bridge를 제공한다.

## 제품 원칙

- **초보자 우선**: blank canvas만 주지 않고 guided mission, component card, glossary, build log를 제공한다.
- **현실 기반 단순화**: 완전한 전기/물리 해석이 아니라, 초보자가 이해할 수 있는 현실적 규칙부터 모델링한다.
- **안전한 실패**: 실제 부품을 태우거나 공간을 차지하지 않고 마음껏 실험할 수 있어야 한다.
- **현실 제작으로 이어지는 데이터**: MVP부터 부품 id, pin, net, dimensions, mounting point, material hint를 잃지 않는 데이터 모델을 둔다.

## MVP 기능 목록

| # | 기능명 | 한 줄 설명 | MVP에 필요한 이유 |
|---|---|---|---|
| 1 | 서버 권한 복셀 월드 | Go 서버가 고정 틱으로 복셀 배치/삭제와 월드 스냅샷을 관리 | 서버 권한 구조 검증 없이는 멀티플레이와 보안 전제가 성립하지 않음 |
| 2 | Three.js 3D 조립 클라이언트 | 클라이언트가 서버 스냅샷을 받아 InstancedMesh로 렌더링하고 블록 편집 이벤트를 전송 | 사용자가 샌드박스의 핵심 조립 경험을 확인해야 함 |
| 3 | 스냅샷 버퍼와 보간 렌더링 | 클라이언트가 서버 tick snapshot 사이를 일정 지연으로 보간해 jitter를 줄임 | 서버 권한 물리 게임의 가장 큰 시각화 리스크를 초기에 검증해야 함 |
| 4 | 최소 회로 엔진 | 전원, 전선, 버튼, AND 게이트, MCU 출력 핀의 High/Low 상태를 서버에서 계산 | 하드웨어-물리 하이브리드의 전자계 절반을 구성 |
| 5 | 최소 액추에이터 연동 | MCU 출력 또는 버튼 신호가 피스톤/모터 블록의 움직임으로 변환 | 회로 신호가 물리 변화로 이어지는 USP를 증명 |
| 6 | 초보자 미션과 컴포넌트 카드 | LED/버튼/모터/피스톤 같은 기본 부품의 역할, 배선법, 현실 주의사항을 앱 안에서 설명 | 사용자가 전문 지식 없이도 첫 성공 경험을 얻어야 함 |

## 유저 스토리 & 수용 기준

### 기능 1: 서버 권한 복셀 월드

**유저 스토리**
> As a builder, I want to place and remove blocks through the server, so that the world state is consistent for every connected client.

**수용 기준**
- Given 서버가 실행 중이고 클라이언트가 연결된 상태에서, When 사용자가 빈 좌표에 블록을 배치하면, Then 서버가 좌표/블록 타입을 검증한 뒤 다음 snapshot에 해당 블록을 포함한다.
- Given 같은 좌표에 두 클라이언트가 동시에 블록을 배치하려 할 때, When 서버가 두 command를 같은 tick에서 처리하면, Then 결정적인 충돌 처리 규칙에 따라 하나의 결과만 snapshot에 반영된다.

### 기능 2: Three.js 3D 조립 클라이언트

**유저 스토리**
> As a builder, I want to see and edit the voxel world in 3D, so that I can build a small machine interactively.

**수용 기준**
- Given 서버에서 초기 world snapshot을 받았을 때, When 클라이언트가 렌더링하면, Then 각 블록 타입이 InstancedMesh 기반으로 화면에 표시된다.
- Given 사용자가 블록 면을 클릭했을 때, When 배치 또는 삭제 모드를 실행하면, Then 클라이언트는 직접 상태를 바꾸지 않고 command를 서버로 전송한다.

### 기능 3: 스냅샷 버퍼와 보간 렌더링

**유저 스토리**
> As a player, I want moving blocks to look smooth even when the server ticks at 20Hz, so that the server-authoritative model still feels responsive.

**수용 기준**
- Given 서버가 tick id, server time, entity transform을 포함한 snapshot을 20Hz로 전송할 때, When 클라이언트가 60FPS 이상으로 렌더링하면, Then render timestamp 기준으로 과거 100-150ms 구간의 두 snapshot 사이를 보간한다.
- Given snapshot 간격이 일시적으로 벌어졌을 때, When 보간 대상 snapshot이 부족하면, Then 클라이언트는 제한된 extrapolation 또는 last-known transform 유지로 튐을 최소화한다.
- Given 개발자가 debug overlay를 켰을 때, When world가 렌더링되면, Then snapshot buffer length, interpolation delay, RTT, dropped snapshot count를 확인할 수 있다.

### 기능 4: 최소 회로 엔진

**유저 스토리**
> As a circuit builder, I want to connect simple signal blocks, so that I can verify logical behavior before driving a physical block.

**수용 기준**
- Given 전원, 전선, AND 게이트, 버튼 블록이 배치된 상태에서, When 서버 tick이 진행되면, Then 회로 그래프가 결정적으로 평가되어 각 노드의 High/Low 상태가 snapshot에 포함된다.
- Given 버튼 입력이 바뀌었을 때, When 다음 tick이 계산되면, Then 연결된 전선과 게이트 출력 상태가 갱신된다.

### 기능 5: 최소 액추에이터 연동

**유저 스토리**
> As a maker, I want a signal output to move a motor or piston block, so that I can see code/circuit output become physical motion.

**수용 기준**
- Given MCU 출력 핀 또는 버튼 신호가 피스톤 블록에 연결된 상태에서, When 신호가 High가 되면, Then 서버는 피스톤의 target transform을 계산하고 snapshot에 반영한다.
- Given 신호가 Low가 되었을 때, When 다음 tick들이 진행되면, Then 피스톤은 서버 계산에 따라 원위치로 돌아간다.

### 기능 6: 초보자 미션과 컴포넌트 카드

**유저 스토리**
> As a beginner maker, I want the app to explain what each part does and what real-world rule it represents, so that I can learn while building instead of needing electronics knowledge upfront.

**수용 기준**
- Given 사용자가 LED, 버튼, PWM pin, motor/driver, piston 같은 부품을 선택했을 때, When inspect panel을 열면, Then 역할, 연결 방법, 실제 주의사항, 시뮬레이션 단순화 범위를 볼 수 있다.
- Given 사용자가 첫 guided mission을 시작했을 때, When LED 켜기, 버튼 입력, PWM 출력, 피스톤 구동을 순서대로 완료하면, Then build log에 완료 단계가 기록된다.
- Given 사용자가 GPIO pin에 motor를 직접 연결하려 할 때, When 서버가 회로를 평가하면, Then invalid wiring 또는 warning을 표시하고 motor driver/transistor component 사용을 안내한다.

## Out of Scope (MVP에서 명시적으로 뺀 것)

- **사용자 계정/권한/소셜 로그인** — MVP는 로컬 또는 LAN 테스트 중심. 서버 권한 구조 검증이 먼저다.
- **영구 DB 기반 월드 저장** — 초기에는 메모리 월드 + JSON save/load 정도로 충분하다.
- **완전한 강체 물리 엔진** — 첫 MVP는 kinematic actuator와 단순 충돌부터 시작한다.
- **대규모 public multiplayer** — 2-4명 개발 테스트 방 기준으로 제한한다.
- **AI 코딩 어시스턴트** — 핵심 샌드박스 루프 검증 후 Later로 이동한다.
- **완전한 Arduino 호환 런타임** — 처음에는 제한된 DSL 또는 Blockly-style command model로 시작한다.
- **모바일 지원** — 데스크톱 브라우저 우선.
- **완성형 3D 프린터 출력/제작 패키지** — MVP에서는 Reality Bridge 데이터 모델만 준비하고, `3MF/STL + BOM + wiring guide + Arduino sketch` 패키지는 v0.2 이후로 이동한다.

## 의존성

- **외부 API**: MVP 없음.
- **인증/권한**: MVP 없음. 로컬 개발 서버와 임시 room id만 사용.
- **플랫폼 요구사항**: 데스크톱 Chromium 계열 브라우저 우선, WebGL2 필요.
- **개발 환경**: Go 1.22+, Node.js 20+, pnpm 또는 npm, Vite, Three.js.

## 비기능 요구사항

- **성능**: 서버 tick 20Hz 기준 32x32x16 테스트 월드와 100개 이하 dynamic entity에서 안정 동작.
- **네트워크**: WebSocket snapshot은 tick id와 server time을 포함하고, client command는 idempotency key를 포함한다.
- **보안/치트 방지**: 클라이언트는 command만 전송하고 authoritative state를 직접 쓰지 않는다.
- **테스트**: Go 서버의 world/circuit/snapshot logic은 단위 테스트 우선 작성. 클라이언트 interpolation은 순수 함수 테스트와 Playwright smoke test로 검증.
- **관측성**: debug overlay와 서버 로그로 tick duration, command queue length, connected clients, snapshot size를 확인한다.

## 가장 큰 리스크

1. 서버 권한 tick rate와 클라이언트 FPS 사이의 간극으로 움직임이 튀어 보일 수 있다.
2. 회로 그래프 평가와 물리 업데이트 순서가 불명확하면 결과가 비결정적으로 보일 수 있다.
3. MVP 범위가 “게임 전체”로 커지면 핵심 기술 검증 전에 프로젝트가 무거워진다.
4. 현실 제작 가능성을 과장하면 신뢰를 잃는다. 시뮬레이션에서 단순화한 부분을 component card에 명확히 드러내야 한다.
5. 사용자의 전문 지식 부족을 앱 밖의 문제로 밀어두면 첫 성공 경험을 만들기 어렵다.

## 검증 목표

- 2명의 클라이언트가 같은 방에 접속해 동일한 복셀 월드를 본다.
- 한 사용자가 회로를 구성하고 버튼을 누르면 다른 사용자 화면에서도 같은 tick 결과가 보인다.
- 서버 20Hz, 클라이언트 60FPS 환경에서 피스톤 또는 모터 움직임이 눈에 띄는 jitter 없이 렌더링된다.
- 초보자가 guided mission을 따라 LED/버튼/피스톤의 관계를 이해하고, 잘못된 motor wiring warning을 볼 수 있다.

## 참고 조사

- Research Notes: [`wire-craft-research-notes.md`](./wire-craft-research-notes.md)
