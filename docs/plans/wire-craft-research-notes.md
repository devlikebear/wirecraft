# Wire Craft — Research Notes & Reality Bridge

_작성일: 2026-04-18_

## 조사 목적

Wire Craft의 방향을 “회로/로봇 게임”이 아니라 “공간, 돈, 장비 제약 없이 메이커 경험을 시작하고, 나중에는 실제 제작까지 이어지는 온라인 공작실”로 잡기 위한 자료 조사다. 사용자의 전문 지식 부족은 약점으로만 볼 게 아니라, 제품이 풀어야 할 핵심 문제로 본다.

## 개인적 출발점

어릴 때 라디오 만들기, 과학상자, 만능회로키트, 전구 켜기, 모터 돌리기 같은 활동을 해보고 싶었지만 경제적 제약 때문에 못 했던 경험이 있다. 어른이 된 뒤에도 실제 로봇을 만들려면 공간, 부품, 장비, 안전, 전문 지식 장벽이 크다. Wire Craft는 이 장벽을 온라인에서 낮추고, “가상에서 마음껏 만들고 실패해본 다음, 원하면 현실 제작으로 옮겨갈 수 있는” 경험을 목표로 한다.

## 경쟁/참고 제품 조사

### Tinkercad Circuits

Autodesk Tinkercad의 가이드에 따르면 Tinkercad Circuits는 전자 부품을 배치/배선하고, Arduino 또는 micro:bit 기반 동작을 block coding 또는 text code로 만들며, schematic view와 시뮬레이션을 제공한다. 특히 실제 회로를 배선하기 전에 가상으로 반응을 확인하는 접근은 Wire Craft의 교육 UX에 직접 참고할 만하다.

**시사점**

- 초보자에게는 3D 월드보다 먼저 “왜 이 부품을 이렇게 연결하는가”를 설명하는 scaffolding이 필요하다.
- Wire Craft는 단순 회로 시뮬레이터를 반복하면 안 된다. 회로 결과가 3D 물리 장치의 움직임으로 이어지는 지점이 핵심 차별화다.
- schematic view는 MVP에는 과하지만, Phase 5 이후 회로를 이해시키는 보조 뷰로 가치가 있다.

출처: [Tinkercad Getting Started Guide](https://images.tinkercad.com/jl5ii4oqrdmc/4nMtGEJi53o3n7V2KJKG7a/0b1818fdecff4868bebc74ebeaf925cf/Tinkercad_Getting_Started_Guide_printable.pdf)

### Wokwi

Wokwi는 Arduino, ESP32, STM32, Pi Pico 등 여러 보드와 센서/부품을 브라우저에서 시뮬레이션하는 온라인 전자 시뮬레이터다. 브라우저만으로 시작할 수 있고, 가상 하드웨어라 망가뜨릴 걱정이 없으며, 디버거/logic analyzer/custom chip 같은 고급 기능도 제공한다. 다만 Wokwi 문서의 potentiometer reference는 full analog simulation을 지원하지 않는다고 명시한다.

**시사점**

- Wire Craft도 초기에 “완전한 아날로그 전기/기계 해석”을 목표로 하면 안 된다.
- 첫 모델은 digital High/Low, PWM duty cycle, kinematic actuator, 단순 충돌처럼 설명 가능한 수준으로 제한한다.
- 대신 사용자가 무엇이 단순화되었는지 알 수 있도록 component card와 warning을 제공한다.

출처: [Wokwi Docs](https://docs.wokwi.com/), [Wokwi potentiometer reference](https://docs.wokwi.com/parts/wokwi-potentiometer)

### KiCad

KiCad는 schematic capture, PCB layout, SPICE simulator, electrical rules checker, 3D viewer를 제공하는 오픈소스 EDA 도구다. 초보자용 샌드박스라기보다는 실제 회로/PCB 설계 도구에 가깝다.

**시사점**

- Wire Craft 내부에 KiCad 수준의 EDA를 재구현하지 않는다.
- 장기적으로는 `export to KiCad-friendly schematic/netlist` 또는 `BOM + wiring guide`를 현실 제작 경로로 제공한다.
- MVP에서는 KiCad export가 아니라 “나중에 KiCad로 옮길 수 있게 부품 id, pin, net 정보를 잃지 않는 데이터 모델”을 먼저 설계한다.

출처: [KiCad](https://www.kicad.org/)

## 초보자 온보딩 조사

### Arduino Education Starter Kit

Arduino Education Starter Kit는 중학생을 대상으로 programming, coding, electronics를 단계적으로 가르치며, 사전 지식이 없어도 따라갈 수 있도록 teacher guide, lesson, engineering logbook을 제공한다. 9개의 90분 수업과 2개의 open-ended group project 구조도 참고할 만하다.

**시사점**

- Wire Craft는 sandbox만 주면 안 된다. “LED 켜기 → 버튼으로 LED 제어 → PWM으로 밝기 조절 → 모터/서보 움직이기 → 센서로 움직임 제어” 같은 미션 경로가 있어야 한다.
- 사용자가 전문 지식이 없어도 진행하도록 component card, glossary, build log, guided mission을 MVP 범위에 최소한 포함한다.
- 교육용/유료화 가능성은 “교사용 수업/과제/공유 room”에서 생긴다. 단, MVP에서는 계정/수업 관리보다 핵심 제작 루프 검증이 먼저다.

출처: [Arduino Education Starter Kit](https://store.arduino.cc/products/arduino-education-starter-kit)

### Blockly

Blockly는 block이 프로그래밍 언어의 expression/statement를 나타내고, block definition과 code generator를 통해 텍스트 코드로 변환되는 구조를 제공한다.

**시사점**

- 초보자 모드는 Blockly 또는 Blockly-like DSL이 적합하다.
- 서버는 사용자가 작성한 코드를 그대로 실행하지 않고, block/code를 제한된 intermediate representation으로 변환해 MCU pin state만 제어한다.
- 텍스트 코딩은 나중에 추가하되, MVP의 핵심은 “핀을 읽고 쓰는 행동 모델”이다.

출처: [Blockly custom blocks](https://developers.google.com/blockly/guides/get-started/blocks)

## 하드웨어 현실성 기준

### PWM/GPIO 모델

Arduino 문서에 따르면 PWM은 `analogWrite()`로 사용하며 기본적으로 8-bit 해상도, 즉 0-255 범위 값을 사용한다. 보드마다 PWM 지원 pin이 다르다.

**Wire Craft 모델링 결정**

- `digitalWrite(pin, HIGH/LOW)`는 boolean output으로 모델링한다.
- `analogWrite(pin, value)`는 실제 아날로그 전압이 아니라 PWM duty cycle로 모델링한다.
- board profile에 pin capability를 둔다. 예: Arduino Uno profile은 PWM pin을 3, 5, 6, 9, 10, 11로 시작한다.

출처: [Arduino Help Center: PWM output](https://support.arduino.cc/hc/en-us/articles/9350537961500-Use-PWM-output-with-Arduino)

### 모터/액추에이터 모델

초보자에게 흔한 위험은 MCU pin에 모터를 직접 연결해도 된다고 오해하는 것이다. Adafruit의 Arduino DC motor lesson은 작은 DC motor도 Arduino digital output이 직접 감당하기 어렵고, transistor와 diode를 사용해 더 큰 motor current와 역전압을 다뤄야 한다고 설명한다.

**Wire Craft 모델링 결정**

- MVP에서 “GPIO pin -> DC motor 직접 연결”은 warning 또는 invalid wiring으로 처리한다.
- motor block은 `motor_driver` 또는 `transistor_switch` component를 통해서만 안전하게 동작하게 한다.
- Phase 3의 피스톤/모터는 물리적으로는 kinematic simulation이지만, 전기적으로는 “driver가 있어야 움직인다”는 규칙을 UI에 노출한다.

출처: [Adafruit Arduino Lesson 13: DC Motors](https://learn.adafruit.com/adafruit-arduino-lesson-13-dc-motors/transistors)

## 현실 제작 연계 조사

### glTF/GLB

Khronos는 glTF를 3D scene/model의 효율적 전송과 로딩을 위한 royalty-free specification으로 설명한다. Three.js는 `GLTFExporter`를 통해 scene을 `.gltf` 또는 `.glb`로 export할 수 있다.

**시사점**

- Wire Craft의 3D 월드 공유/미리보기/웹 임베드는 glTF/GLB가 자연스럽다.
- glTF는 runtime 3D asset에 적합하지만, 3D 프린팅 제조 지시 전체를 담는 포맷은 아니다.

출처: [Khronos glTF](https://www.khronos.org/gltf/), [Three.js GLTFExporter](https://threejs.org/docs/pages/GLTFExporter.html)

### 3MF/STL/OBJ

3MF Consortium은 3MF를 additive manufacturing을 위한 full-fidelity 3D model 전송 포맷으로 설명한다. Autodesk Fusion은 mesh body를 3MF, STL, OBJ로 export할 수 있고, STL은 모델을 flat facets로 근사한 stereolithography용 형식으로 설명한다.

**시사점**

- 실제 3D 프린터 연계는 장기적으로 3MF를 우선 검토하고, STL은 범용 fallback으로 둔다.
- 그러나 print file만으로는 로봇이 완성되지 않는다. 전자 부품 BOM, wiring guide, Arduino sketch, 조립 순서, 안전 warning이 같이 필요하다.
- 따라서 Wire Craft의 장기 산출물은 단일 STL이 아니라 `Reality Pack`이어야 한다.

출처: [3MF Specification](https://3mf.io/spec/), [Autodesk Fusion mesh export](https://help.autodesk.com/cloudhelp/ENU/Fusion-Mesh/files/MESH-EXPORT-TOOLS.htm), [Autodesk Fusion export designs](https://help.autodesk.com/cloudhelp/ENU/Fusion-Assemble/files/ASM-EXPORT-DESIGN.htm)

## 제품 결정

### Primary Persona

“어릴 때 메이커 키트와 로봇 제작을 해보고 싶었지만 돈, 공간, 장비, 지식 장벽 때문에 못 했고, 이제 온라인에서 안전하게 실험하며 실제 제작으로 넘어가고 싶은 성인 초보 메이커.”

### Secondary Personas

- 중고등학생에게 회로/로봇/코딩을 가르치고 싶은 교사
- 3D 프린터와 아두이노를 갖고 있지만 설계-배선-코드 연결이 막히는 취미 제작자
- 온라인으로 함께 공장 라인/로봇/장치를 만들고 싶은 소규모 팀

### Product North Star

사용자가 “가상 공간에서 직접 배선하고 코딩한 장치가 움직이는 걸 본 뒤, 필요한 부품과 출력 파일을 받아 현실에서도 다시 만들어볼 수 있다”고 느끼게 한다.

### MVP 원칙

- 완벽한 전기/물리 해석보다 “초보자가 배울 수 있는 현실 기반 단순화”가 우선이다.
- 모든 부품에는 `what it does`, `how to wire`, `real-world warning`, `simulated simplification`을 담은 component card를 둔다.
- 초보자는 blank canvas가 아니라 guided mission에서 시작한다.
- 현실 제작 연계는 MVP에서 완성하지 않고, 데이터 모델에 부품 id, pin, net, mounting point, dimensions, material hint를 남겨 v0.2의 `Reality Pack`으로 이어지게 한다.

## Reality Pack 방향

v0.2 이후 목표 산출물:

- `model.glb`: 웹 공유/미리보기용 3D 모델
- `print.3mf` 또는 `print.stl`: 출력 가능한 구조물 후보
- `bom.csv`: 부품 목록, 수량, 대체품, 예상 가격
- `wiring.md` 또는 `wiring.json`: 핀 연결표와 배선 순서
- `controller.ino`: Arduino-style starter sketch
- `constraints.json`: 전압, 전류, torque, clearance, material, print tolerance 같은 제한 메타데이터

MVP에서는 이 중 `bom.csv`, `wiring.md`, `model.glb`의 prototype export를 Phase 5의 optional stretch로 둔다.

