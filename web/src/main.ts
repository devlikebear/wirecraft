import './styles.css';
import {
  AmbientLight,
  BoxGeometry,
  Color,
  DirectionalLight,
  GridHelper,
  Mesh,
  MeshStandardMaterial,
  PerspectiveCamera,
  Scene,
  SRGBColorSpace,
  Vector2,
  WebGLRenderer,
} from 'three';

const app = document.querySelector<HTMLElement>('#app');

if (!app) {
  throw new Error('Missing #app root element');
}

const scene = new Scene();
scene.background = new Color(0x101211);

const camera = new PerspectiveCamera(50, window.innerWidth / window.innerHeight, 0.1, 100);
camera.position.set(9, 8, 10);
camera.lookAt(0, 0, 0);

const renderer = new WebGLRenderer({ antialias: true, preserveDrawingBuffer: true });
renderer.outputColorSpace = SRGBColorSpace;
renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
renderer.setSize(window.innerWidth, window.innerHeight);
app.appendChild(renderer.domElement);

const ambient = new AmbientLight(0xf1ead9, 0.45);
scene.add(ambient);

const keyLight = new DirectionalLight(0xffffff, 1.8);
keyLight.position.set(5, 8, 4);
scene.add(keyLight);

const fillLight = new DirectionalLight(0xa7c5bd, 0.7);
fillLight.position.set(-6, 4, -5);
scene.add(fillLight);

const grid = new GridHelper(18, 18, 0x53635d, 0x27302d);
scene.add(grid);

const blockGeometry = new BoxGeometry(1, 1, 1);
const baseMaterial = new MeshStandardMaterial({ color: 0x6d8f7d, roughness: 0.7 });
const copperMaterial = new MeshStandardMaterial({ color: 0xb56a3c, roughness: 0.58, metalness: 0.15 });
const moverMaterial = new MeshStandardMaterial({ color: 0xd6c064, roughness: 0.45 });

const baseBlocks: Mesh[] = [];
const footprint = [
  [-2, 0, -1],
  [-1, 0, -1],
  [0, 0, -1],
  [1, 0, -1],
  [-2, 0, 0],
  [1, 0, 0],
  [-2, 0, 1],
  [-1, 0, 1],
  [0, 0, 1],
  [1, 0, 1],
];

for (const [x, y, z] of footprint) {
  const block = new Mesh(blockGeometry, baseMaterial);
  block.position.set(x, y + 0.5, z);
  scene.add(block);
  baseBlocks.push(block);
}

const copperTrace = new Mesh(new BoxGeometry(4.2, 0.08, 0.16), copperMaterial);
copperTrace.position.set(-0.5, 1.04, 0);
scene.add(copperTrace);

const debugMover = new Mesh(blockGeometry, moverMaterial);
debugMover.position.set(0, 1.6, 0);
scene.add(debugMover);

const pointer = new Vector2();
let targetYaw = 0;

window.addEventListener('pointermove', (event) => {
  pointer.x = event.clientX / window.innerWidth - 0.5;
  pointer.y = event.clientY / window.innerHeight - 0.5;
  targetYaw = pointer.x * 0.25;
});

window.addEventListener('resize', () => {
  camera.aspect = window.innerWidth / window.innerHeight;
  camera.updateProjectionMatrix();
  renderer.setSize(window.innerWidth, window.innerHeight);
});

function animate(timeMs: number) {
  const time = timeMs / 1000;
  debugMover.position.y = 1.6 + Math.sin(time * 2.2) * 0.24;
  debugMover.rotation.y = time * 0.7;
  scene.rotation.y += (targetYaw - scene.rotation.y) * 0.035;

  renderer.render(scene, camera);
  requestAnimationFrame(animate);
}

requestAnimationFrame(animate);
