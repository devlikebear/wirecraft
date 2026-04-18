# Wire Craft — 개발 로드맵

_이 문서는 전체 개발 계획의 총괄본이다. Claude Code는 이 로드맵을 먼저 읽고 전체 맥락을 파악한 뒤, 각 페이즈 작업지시서로 이동해서 실제 구현을 진행한다. 페이즈 간 전환 시 반드시 사용자 확인을 받는다._

_작성일: 2026-04-18_
_예상 전체 소요: 6-10주_
_페이즈 수: 5개_

## Overview

Wire Craft는 서버 권한 구조를 기반으로 복셀 조립, 회로 시뮬레이션, 물리 액추에이터를 하나의 협업형 샌드박스에 결합하는 프로젝트다. 핵심 차별점은 회로 신호가 실제 3D 월드의 움직임을 만들고, 그 상태를 서버가 결정적으로 계산해 모든 클라이언트에 동기화한다는 점이다.

가장 먼저 검증해야 할 것은 콘텐츠 양이 아니라 네트워크/시뮬레이션 루프다. 따라서 Phase 1부터 서버 tick snapshot, client command, interpolation buffer를 포함한다. 보간은 Phase 4의 최적화가 아니라 MVP 성립 조건이다.

제품 방향은 “온라인에서 마음껏 로봇과 회로를 만드는 공작실”이다. 사용자가 전문 지식 없이도 guided mission과 component card로 시작하고, 장기적으로는 가상 제작물을 부품 목록, 배선표, 코드, 3D 출력 후보 파일로 옮기는 Reality Bridge를 제공한다.

## 완료 조건 (전체)

- [ ] Go 서버가 authoritative world, circuit, actuator state를 고정 tick으로 계산한다.
- [ ] Three.js 클라이언트가 서버 snapshot을 받아 InstancedMesh 기반 복셀 월드를 렌더링한다.
- [ ] 클라이언트가 snapshot buffer 기반 interpolation으로 dynamic entity를 부드럽게 표시한다.
- [x] 최소 회로(전원, 버튼, 전선, 게이트, MCU output)가 서버에서 결정적으로 평가된다.
- [ ] 회로 신호가 피스톤 또는 모터의 물리적 움직임으로 연결된다.
- [ ] 2-4명의 클라이언트가 같은 room에 접속해 동일한 상태를 본다.
- [ ] 초보자용 component card와 guided mission이 최소 1개 동작한다.
- [ ] Reality Bridge를 위한 부품 id, pin/net, dimensions, mounting point 메타데이터가 데이터 모델에 남는다.
- [ ] 각 페이즈별 테스트와 수동 체크포인트가 통과한다.

## 기술 스택 / 환경

- **서버**: Go, 표준 `net/http`, WebSocket 라이브러리(`nhooyr.io/websocket` 또는 `github.com/coder/websocket` 중 구현 시 선택)
- **클라이언트**: Vite, TypeScript, Three.js
- **상태 동기화**: WebSocket command/snapshot protocol, tick id, server time, sequence id
- **물리**: 초기 kinematic simulation 직접 구현. 추후 Go physics library 검토
- **저장소 / DB**: MVP는 메모리 + JSON save/load. 영구 DB는 Later
- **현실 제작 연계**: MVP는 metadata/BOM 기반 준비. v0.2 이후 glTF/GLB, 3MF/STL, wiring guide, Arduino sketch export 검토
- **테스트**: Go `testing`, TypeScript unit test(Vitest 권장), Playwright smoke test
- **실행 환경**: 로컬 개발 서버, 데스크톱 Chromium 우선

## Out of Scope

- Public account system, billing, moderation
- 완전한 Arduino/C++ 런타임 호환
- 대규모 MMO 수준의 interest management
- 완전한 CAD/EDA 도구 수준의 회로 해석
- 완성형 3D 프린터 제조 패키지
- 모바일/터치 UI
- AI code assistant

## 규모 판단

기능은 5개 이상이고, 서버 tick loop, 네트워크 프로토콜, Three.js 렌더링, 회로 엔진, 물리 동기화가 모두 포함된다. 예상 구현 시간이 16시간을 크게 넘고 각 단계마다 검증 포인트가 명확하므로 단일 `dev-plan.md`가 아니라 `roadmap.md + phase-N-*.md` 구조가 맞다.

## 페이즈 구성

각 페이즈는 수직 슬라이스(vertical slice)다. 끝나면 동작하는 무언가가 있어야 한다.

### Phase 1: Authoritative Voxel Loop — 서버 권한 3D 편집 루프

Status: Completed. Tracking issue: [#1](https://github.com/devlikebear/wirecraft/issues/1).

- **목표**: 한 명 이상의 클라이언트가 서버에 접속해 복셀을 배치/삭제하고, 서버 snapshot을 보간 렌더링하는 최소 샌드박스가 동작한다.
- **포함 기능**: Go 서버 tick loop, WebSocket protocol, voxel world model, Three.js rendering, raycast edit command, interpolation buffer
- **예상 소요**: 1-2주
- **작업지시서**: [`wire-craft-phase-1-authoritative-voxel-loop.md`](./wire-craft-phase-1-authoritative-voxel-loop.md)
- **Checkpoint 요약**: 2개 브라우저 창에서 같은 월드를 보고, dynamic test entity가 jitter 없이 움직인다.

### Phase 2: Circuit Runtime — 서버 사이드 회로 엔진

Status: Completed. Tracking issue: [#17](https://github.com/devlikebear/wirecraft/issues/17).

- **목표**: 서버가 전원, 버튼, 전선, 게이트, MCU output node의 High/Low 상태를 tick마다 결정적으로 계산하고, 초보자가 각 부품의 의미를 component card로 확인한다.
- **포함 기능**: circuit graph model, signal propagation, block-to-circuit mapping, debug visualization, component card
- **예상 소요**: 1-2주
- **작업지시서**: [`wire-craft-phase-2-circuit-runtime.md`](./wire-craft-phase-2-circuit-runtime.md)
- **Checkpoint 요약**: 버튼과 게이트 조합 결과가 모든 클라이언트에서 같은 tick에 표시된다.

### Phase 3: Physical Actuators — 회로 신호와 물리 블록 결합

Status: Verification complete; awaiting user approval. Tracking issue: [#28](https://github.com/devlikebear/wirecraft/issues/28).

- **목표**: 회로 출력이 피스톤 또는 모터 블록 움직임으로 변환되고, 클라이언트는 transform snapshot을 보간 렌더링한다.
- **포함 기능**: actuator component, simple kinematic movement, sensor input stub, transform replication
- **예상 소요**: 1-2주
- **작업지시서**: [`wire-craft-phase-3-physical-actuators.md`](./wire-craft-phase-3-physical-actuators.md)
- **Checkpoint 요약**: 사용자가 버튼을 누르면 피스톤이 서버 계산으로 움직이고 다른 클라이언트에도 동일하게 보인다.

### Phase 4: Multiplayer Physics Sync — 협업과 동기화 강화

- **목표**: 2-4명 협업 환경에서 동시 편집, command conflict, snapshot delta, tick performance를 안정화한다.
- **포함 기능**: room/session model, command queue, conflict resolution, delta snapshot, basic collision constraints
- **예상 소요**: 1-2주
- **작업지시서**: [`wire-craft-phase-4-multiplayer-physics-sync.md`](./wire-craft-phase-4-multiplayer-physics-sync.md)
- **Checkpoint 요약**: 2-4개 클라이언트가 동시에 편집해도 서버 상태가 깨지지 않고 snapshot size와 tick time이 관측된다.

### Phase 5: Blueprint & Reality Bridge Prep — 재사용성과 현실 제작 준비

- **목표**: 작은 회로/장치 영역을 blueprint로 저장/불러오고, Reality Pack으로 확장 가능한 BOM/wiring/model metadata를 준비한다.
- **포함 기능**: blueprint JSON, local save/load, toolbar/mode UI, circuit overlay, demo scenario, BOM/wiring metadata prototype
- **예상 소요**: 1-2주
- **작업지시서**: [`wire-craft-phase-5-blueprint-ux-polish.md`](./wire-craft-phase-5-blueprint-ux-polish.md)
- **Checkpoint 요약**: 사용자가 만든 버튼-게이트-피스톤 장치를 blueprint로 저장하고 새 위치에 배치할 수 있다.

## 페이즈 간 의존성

```text
Phase 1: authoritative voxel + snapshot/interpolation foundation
  -> Phase 2: circuit state on authoritative server
    -> Phase 3: circuit output drives physical actuator
      -> Phase 4: multiplayer conflict/sync/performance hardening
        -> Phase 5: blueprint and Reality Bridge prep
```

**병렬 가능성**: Phase 2의 circuit engine 순수 로직과 Phase 1의 Three.js 렌더링 세부 구현은 일부 병렬 가능하다. 다만 WebSocket protocol과 snapshot schema는 Phase 1에서 먼저 고정해야 한다.

## 페이즈 간 전환 규칙

각 페이즈는 다음 조건을 만족해야 완료로 간주한다.

1. 작업지시서의 모든 체크박스 완료
2. Checkpoint 블록의 자동/수동 검증 통과
3. 사용자가 명시적으로 "Phase N 완료 확인, 다음 진행" 승인

Checkpoint가 실패하면 Claude Code는 실패 항목, 재현 방법, 추정 원인을 보고한 뒤 수정하고 재검증한다.

## 최종 완료 체크리스트

- [ ] 모든 페이즈 Checkpoint 통과
- [ ] 서버 테스트: `go test ./...`
- [ ] 클라이언트 테스트: `npm test`
- [ ] 클라이언트 빌드: `npm run build`
- [ ] 브라우저 smoke test: Playwright로 world 접속, block place/remove, button-actuator demo 확인
- [ ] 2개 브라우저 창에서 같은 room 접속 후 동기화 확인
- [ ] Out of Scope 항목 미구현 확인
- [ ] README에 실행 방법과 MVP 범위 기록

## 핵심 설계 메모

### Tick Rate와 FPS 간극 처리

서버는 20Hz 또는 30Hz 고정 tick으로 authoritative snapshot을 발행한다. 클라이언트는 최신 snapshot을 즉시 그리지 않고 100-150ms 정도 늦은 render timestamp를 기준으로 과거 snapshot 두 개를 찾아 보간한다. 이렇게 하면 네트워크 지터가 있어도 대부분의 프레임에서 보간 대상이 존재한다.

기본 알고리즘:

1. snapshot 수신 시 `snapshotBuffer`에 `tick`, `serverTimeMs`, `entities`, `blocksChanged`를 저장한다.
2. 렌더 루프에서 `renderServerTime = estimatedServerNow - interpolationDelayMs`를 계산한다.
3. buffer에서 `before.serverTimeMs <= renderServerTime <= after.serverTimeMs`인 두 snapshot을 찾는다.
4. static voxel은 delta를 즉시 반영하고, dynamic entity transform은 `alpha`로 position/quaternion/scale을 보간한다.
5. `after`가 없으면 짧은 시간만 extrapolate하고, 한계를 넘으면 last-known transform을 유지한다.

이 기능은 Phase 1에서 구현한다. 나중에 붙이면 서버 권한 구조가 실제로 플레이 가능한지 늦게 알게 된다.

### Reality Bridge 전략

실제 제작 연계는 “STL 하나 뽑기”가 아니다. 로봇/회로 제작에는 출력 가능한 구조물, 부품 목록, 배선표, 컨트롤러 코드, 전기적 주의사항이 함께 필요하다. 따라서 Wire Craft의 장기 export 단위는 `Reality Pack`으로 정의한다.

Reality Pack 후보:

- `model.glb`: 웹 공유/미리보기용 3D 모델
- `print.3mf` 또는 `print.stl`: 3D 프린터 출력 후보
- `bom.csv`: 부품명, 수량, 규격, 대체품, 예상 가격
- `wiring.md` 또는 `wiring.json`: pin/net 연결표
- `controller.ino`: Arduino-style starter sketch
- `constraints.json`: 전압, 전류, torque, clearance, material, print tolerance 제한

MVP에서는 이 전체를 완성하지 않는다. 대신 Phase 2-5에서 부품 id, pin, net, dimensions, mounting point, material hint를 잃지 않도록 데이터 모델을 만든다.

## 참고 자료

- PRD: [`wire-craft-prd.md`](./wire-craft-prd.md)
- Research Notes: [`wire-craft-research-notes.md`](./wire-craft-research-notes.md)
- Phase 1: [`wire-craft-phase-1-authoritative-voxel-loop.md`](./wire-craft-phase-1-authoritative-voxel-loop.md)
- Phase 2: [`wire-craft-phase-2-circuit-runtime.md`](./wire-craft-phase-2-circuit-runtime.md)
- Phase 3: [`wire-craft-phase-3-physical-actuators.md`](./wire-craft-phase-3-physical-actuators.md)
- Phase 4: [`wire-craft-phase-4-multiplayer-physics-sync.md`](./wire-craft-phase-4-multiplayer-physics-sync.md)
- Phase 5: [`wire-craft-phase-5-blueprint-ux-polish.md`](./wire-craft-phase-5-blueprint-ux-polish.md)

---
_Claude Code 사용 시: 이 로드맵을 먼저 읽고 전체 맥락 파악. 그다음 Phase 1 작업지시서로 이동해서 구현 시작. 페이즈 전환은 반드시 사용자 명시적 승인 후._
