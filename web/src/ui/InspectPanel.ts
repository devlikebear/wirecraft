import type { ComponentCard } from '../state/componentCards';

export interface InspectPanel {
  element: HTMLElement;
  setCard(card: ComponentCard | null): void;
}

export function inspectPanelText(card: ComponentCard | null): string {
  if (!card) {
    return 'Select a circuit component.';
  }

  return [
    card.name,
    `Role: ${card.role}`,
    `Pins: ${card.pins.map((pin) => `${pin.name} ${pin.direction}`).join(', ')}`,
    `Wiring: ${card.wiringNotes.join(' ')}`,
    ...card.warnings.map((warning) => `Warning: ${warning}`),
    `Simplification: ${card.simplificationNotes.join(' ')}`,
  ].join('\n');
}

export function createInspectPanel(initialCard: ComponentCard | null): InspectPanel {
  const element = document.createElement('aside');
  element.className = 'inspect-panel';
  element.setAttribute('aria-live', 'polite');

  function setCard(card: ComponentCard | null): void {
    element.replaceChildren();

    if (!card) {
      const empty = document.createElement('p');
      empty.className = 'inspect-panel__empty';
      empty.textContent = inspectPanelText(null);
      element.appendChild(empty);
      return;
    }

    const heading = document.createElement('h2');
    heading.className = 'inspect-panel__title';
    heading.textContent = card.name;
    element.appendChild(heading);

    const role = document.createElement('p');
    role.className = 'inspect-panel__role';
    role.textContent = card.role;
    element.appendChild(role);

    element.appendChild(section('Pins', card.pins.map((pin) => `${pin.name}: ${pin.direction}, ${pin.signal}`)));
    element.appendChild(section('Wiring', card.wiringNotes));
    element.appendChild(section('Warning', card.warnings));
    element.appendChild(section('Simplification', card.simplificationNotes));
  }

  setCard(initialCard);

  return {
    element,
    setCard,
  };
}

function section(title: string, entries: string[]): HTMLElement {
  const container = document.createElement('section');
  container.className = 'inspect-panel__section';

  const heading = document.createElement('h3');
  heading.className = 'inspect-panel__section-title';
  heading.textContent = title;
  container.appendChild(heading);

  const list = document.createElement('ul');
  list.className = 'inspect-panel__list';
  for (const entry of entries) {
    const item = document.createElement('li');
    item.textContent = entry;
    list.appendChild(item);
  }
  container.appendChild(list);

  return container;
}
