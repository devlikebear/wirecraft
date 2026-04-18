# Phase 5: Blueprint & Reality Bridge Prep — 작업지시서

_작성일: 2026-04-18_
_속한 로드맵: [`wire-craft-roadmap.md`](./wire-craft-roadmap.md)_
_예상 소요: 1-2주_

## 페이즈 목표

MVP를 실제로 반복 사용 가능한 데모로 다듬는다. 사용자는 작은 장치 영역을 blueprint로 저장하고 다른 위치에 다시 배치할 수 있으며, 3D 조립 화면과 회로 상태 UI를 오가며 버튼-게이트-액추에이터 demo를 완성할 수 있다. 동시에 v0.2의 Reality Pack으로 이어질 BOM/wiring/model metadata prototype을 준비한다.

## 전제 조건

- [ ] Phase 4 완료 및 사용자 승인
- [ ] block, circuit, actuator state를 직렬화할 수 있는 내부 구조 확보

## 포함 기능

1. blueprint JSON schema
2. selection box
3. save/load
4. placement preview
5. toolbar/mode UX polish
6. demo scenario와 README 정리
7. Reality Pack metadata prototype

## 이 페이즈에서 하지 않는 것

- 서버 DB 저장 → Later
- user-generated marketplace → Later
- AI coding assistant → Later
- 모바일 UI → Out of Scope
- 완성형 3MF/STL 제조 export → v0.2 이후

## 작업 체크리스트

### 작업 그룹 A: Blueprint schema

- [ ] **T5.A.1** — blueprint data model 정의.
  - 파일:
    - `internal/blueprint/blueprint.go`
    - `internal/blueprint/blueprint_test.go`
  - 내용:
    - name, version
    - bounding box
    - blocks with relative positions
    - circuit metadata
    - actuator metadata
  - 테스트:
    - world region -> blueprint 변환
    - blueprint -> world placement 변환
  - 검증: `go test ./internal/blueprint/...`

- [ ] **T5.A.2** — JSON encode/decode 구현.
  - 파일:
    - `internal/blueprint/json.go`
    - `internal/blueprint/json_test.go`
  - 내용:
    - stable JSON field names
    - version field 필수
    - unknown version 처리
  - 검증: `go test ./internal/blueprint/...`

### 작업 그룹 B: Server commands

- [ ] **T5.B.1** — selection export command 구현.
  - 파일:
    - `internal/netproto/command.go`
    - `internal/sim/blueprint_commands.go`
  - 내용:
    - command type: `export_blueprint`
    - selection min/max
    - 서버가 region validation 후 JSON payload 반환
  - 검증: `go test ./...`

- [ ] **T5.B.2** — blueprint placement command 구현.
  - 파일:
    - `internal/sim/blueprint_commands.go`
    - `internal/sim/blueprint_commands_test.go`
  - 내용:
    - command type: `place_blueprint`
    - target origin
    - collision/occupied cell 검증
    - 성공 시 blocks/entities/circuit metadata 생성
  - 검증: `go test ./...`

### 작업 그룹 C: Client UX

- [ ] **T5.C.1** — toolbar와 mode UI 정리.
  - 파일:
    - `web/src/ui/Toolbar.ts`
    - `web/src/styles.css`
  - 내용:
    - modes: build, erase, inspect, select, blueprint
    - icons 또는 짧은 labels 사용
    - active mode 시각 표시
  - 검증: `cd web && npm run build`

- [ ] **T5.C.2** — selection box UI 구현.
  - 파일:
    - `web/src/input/SelectionController.ts`
    - `web/src/render/SelectionRenderer.ts`
  - 내용:
    - 시작 좌표와 끝 좌표 지정
    - bounding box preview
    - selected volume size 표시
  - 검증: 브라우저에서 selection box 표시 확인

- [ ] **T5.C.3** — blueprint save/load UI 구현.
  - 파일:
    - `web/src/ui/BlueprintPanel.ts`
    - `web/src/state/blueprintStore.ts`
  - 내용:
    - export 결과 JSON을 localStorage 또는 file download로 저장
    - JSON import
    - placement preview
  - 검증: 저장한 blueprint를 새 위치에 배치 확인

### 작업 그룹 D: Demo polish

- [ ] **T5.D.1** — demo world seed 추가.
  - 파일:
    - `internal/world/demo.go`
    - `internal/world/demo_test.go`
  - 내용:
    - 버튼 + AND gate + 피스톤 예제
    - 서버 실행 옵션으로 demo world 로드
  - 검증: `go test ./internal/world/...`

- [ ] **T5.D.2** — README와 manual test 문서 정리.
  - 파일:
    - `README.md`
    - `docs/manual-tests/mvp-demo.md`
  - 내용:
    - 서버 실행
    - 클라이언트 실행
    - 2-client multiplayer 확인
    - button-gate-piston demo
    - blueprint save/load
  - 검증: 문서대로 새 터미널에서 실행 가능

### 작업 그룹 E: Reality Bridge prototype

- [ ] **T5.E.1** — BOM export prototype을 구현한다.
  - 파일:
    - `internal/blueprint/bom.go`
    - `internal/blueprint/bom_test.go`
  - 내용:
    - blueprint 안의 component id별 수량 집계
    - 필드: `componentId`, `displayName`, `quantity`, `notes`
    - CSV export 함수
  - 검증: `go test ./internal/blueprint/...`

- [ ] **T5.E.2** — wiring guide export prototype을 구현한다.
  - 파일:
    - `internal/blueprint/wiring.go`
    - `internal/blueprint/wiring_test.go`
  - 내용:
    - circuit graph의 pin/net 연결을 사람이 읽을 수 있는 순서로 변환
    - 예: `Arduino D3 -> resistor -> LED anode`, `LED cathode -> GND`
    - invalid/warning 상태를 함께 출력
  - 검증: `go test ./internal/blueprint/...`

- [ ] **T5.E.3** — 3D model export feasibility note를 작성한다.
  - 파일:
    - `docs/research/reality-pack-export.md`
  - 내용:
    - glTF/GLB: 웹 공유와 미리보기용
    - 3MF: additive manufacturing용 장기 후보
    - STL: 범용 프린팅 fallback이지만 부품/배선/재료 metadata는 별도 필요
    - v0.2에서 실제 export를 구현할 때 필요한 library 후보와 검증 방법 정리
  - 검증: 문서에 export 대상, 한계, 다음 작업이 명시되어 있음

---

## ✅ Phase 5 Checkpoint

**구현 확인:**
- [ ] 사용자가 영역을 선택해 blueprint JSON으로 저장할 수 있다.
- [ ] 저장한 blueprint를 새 위치에 배치할 수 있다.
- [ ] toolbar/mode UI가 build, erase, inspect, select, blueprint 흐름을 지원한다.
- [ ] demo world로 핵심 시나리오를 빠르게 재현할 수 있다.
- [ ] README만 보고 로컬 실행이 가능하다.
- [ ] blueprint에서 BOM과 wiring guide prototype을 만들 수 있다.
- [ ] Reality Pack export 한계와 v0.2 작업 항목이 문서화되어 있다.

**자동 검증:**
- [ ] 서버 테스트 통과: `go test ./...`
- [ ] 클라이언트 테스트 통과: `cd web && npm test`
- [ ] 클라이언트 빌드 통과: `cd web && npm run build`

**수동 확인:**
- [ ] demo world를 열고 버튼-게이트-피스톤 장치가 동작하는지 확인한다.
- [ ] 장치 영역을 blueprint로 저장한다.
- [ ] 빈 위치에 blueprint를 배치하고 동일하게 동작하는지 확인한다.
- [ ] blueprint의 BOM/wiring guide prototype을 확인한다.
- [ ] 브라우저 창 2개에서 blueprint placement 결과가 동일하게 보인다.

**완료 처리:**
1. 위 항목 모두 통과 시 전체 MVP 완료 보고를 작성한다.
2. 사용자가 확인한 뒤 다음 범위(영구 저장, 더 정교한 물리, AI assistant 등)를 새 계획으로 분리한다.
3. 실패 시 실패 항목 보고 → 원인 분석 → 수정 → 재검증한다.

---

## 메모 / 주의

- blueprint는 “재사용 가능한 복잡도 관리”와 “현실 제작으로 이어지는 Reality Bridge”에 직접 연결된다.
- DB 저장이나 공유 marketplace로 확장하고 싶어도 MVP에서는 local save/load와 export prototype까지만 한다.

---
_다음 단계: MVP 완료 후 v0.2 Reality Pack 범위 재기획_
