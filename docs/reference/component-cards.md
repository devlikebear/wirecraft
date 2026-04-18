# WireCraft Component Cards

_Last updated: 2026-04-18_

Component cards describe beginner-facing circuit parts before WireCraft has a full curriculum or search system. Each card keeps a small bridge between the simplified voxel circuit runtime and real-world hardware concerns.

## Schema

- `id`: stable machine identifier.
- `name`: user-facing component name.
- `role`: one-sentence purpose.
- `pins`: named connection points with direction and signal type.
- `wiringNotes`: practical starter wiring hints.
- `warnings`: real hardware caveats that should not be hidden from learners.
- `simplificationNotes`: explicit differences between WireCraft and real electronics.

## Starter Set

| ID | Name | Role |
| --- | --- | --- |
| `power` | Power | Provides a constant High digital signal |
| `wire` | Wire | Carries a digital signal between adjacent components |
| `button` | Button | User-controlled digital switch |
| `and_gate` | AND Gate | Outputs High only when two inputs are High |
| `mcu_output` | MCU Output | Represents a microcontroller output pin or observed endpoint |
| `led` | LED | Visual indicator for a digital output |
| `resistor` | Resistor | Limits current or creates a pull-up/pull-down path |

## Notes By Component

### Power

- Pins: `OUT` output, digital high.
- Wiring: connect `OUT` to a wire, logic input, or indicator input.
- Warning: real supplies need voltage and current limits before powering hardware.
- Simplification: WireCraft treats this as an ideal always-on digital source for starter circuits.

### Wire

- Pins: `A` bidirectional digital, `B` bidirectional digital.
- Wiring: place wires next to circuit blocks to bridge their signals through the voxel grid.
- Warning: real wires have resistance and safe current ratings.
- Simplification: starter wires propagate High, Low, or Unknown without analog voltage drop.

### Button

- Pins: `IN` input digital, `OUT` output digital.
- Wiring: use `set_button` commands to press or release the button in the authoritative simulation.
- Warning: real buttons bounce and often need pull-up or pull-down resistors.
- Simplification: WireCraft models a released button as Low and a pressed button as High.

### AND Gate

- Pins: `A` input digital, `B` input digital, `OUT` output digital.
- Wiring: feed two input paths into the gate, then connect `OUT` to another wire or output block.
- Warning: real logic chips need a matching supply voltage and ground reference.
- Simplification: the current starter runtime infers gate inputs from adjacent circuit blocks.

### MCU Output

- Pins: `IN` input digital.
- Wiring: connect a wire or gate output to observe the resulting server-authoritative signal.
- Warning: real MCU pins cannot drive motors or high-current loads directly.
- Simplification: WireCraft uses this block as a readable digital endpoint before code runtime exists.

### LED

- Pins: `A` input digital high, `K` input ground reference.
- Wiring: drive the anode through a resistor and connect the cathode to ground in real circuits.
- Warning: real LEDs need current-limiting resistors to avoid damage.
- Simplification: LED is documented as a starter card before it becomes a placeable voxel block.

### Resistor

- Pins: `A` passive analog, `B` passive analog.
- Wiring: use resistors with LEDs and buttons in real hardware to control current and default voltage.
- Warning: choose resistance and power rating for the actual voltage and load.
- Simplification: WireCraft documents resistors now but does not simulate analog resistance in Phase 2.
