import './styles.css';
import { DebugOverlay, calculateFps } from './debug/DebugOverlay';
import { EditController } from './input/EditController';
import { BlockType, type Position } from './net/protocol';
import { SnapshotSocket } from './net/socket';
import { EntityRenderer } from './render/EntityRenderer';
import { VoxelRenderer } from './render/VoxelRenderer';
import { DEFAULT_INTERPOLATION_DELAY_MS } from './sim/interpolation';
import { SnapshotStore } from './state/snapshotStore';
import { createToolbar } from './ui/Toolbar';
import {
  AmbientLight,
  Color,
  DirectionalLight,
  GridHelper,
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

app.dataset.interpolationDelayMs = String(DEFAULT_INTERPOLATION_DELAY_MS);

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

const voxelRenderer = new VoxelRenderer(scene);
const entityRenderer = new EntityRenderer(scene);
const debugOverlay = new DebugOverlay(app);
const snapshots = new SnapshotStore({ maxSnapshots: 64 });
const clientId = crypto.randomUUID();
let selectedBlockType = BlockType.Solid;
let wsStatus = 'idle';
let serverTick: number | null = null;
let voxelBlocks = 0;
let latestFps: number | null = null;
let previousFrameTimestampMs: number | null = null;

const snapshotSocket = new SnapshotSocket({
  onSnapshot: (snapshot) => {
    snapshots.append(snapshot);
    voxelRenderer.update(snapshot);
    serverTick = snapshot.tick;
    voxelBlocks = snapshot.blocks.length;
    app.dataset.serverTick = String(snapshot.tick);
    app.dataset.snapshotBuffer = String(snapshots.length);
    app.dataset.voxelBlocks = String(snapshot.blocks.length);
  },
  onStatusChange: (status) => {
    wsStatus = status;
    app.dataset.wsStatus = status;
    console.info(`[wirecraft] websocket ${status}`);
  },
  onError: (error) => {
    console.warn('[wirecraft] websocket error', error);
  },
});
snapshotSocket.connect();

const toolbar = createToolbar({
  selectedBlockType,
  onSelectBlockType: (blockType) => {
    selectedBlockType = blockType;
    app.dataset.selectedBlockType = String(blockType);
  },
});
app.appendChild(toolbar.element);
app.dataset.selectedBlockType = String(selectedBlockType);

const editController = new EditController({
  camera,
  renderer,
  worldRoot: scene,
  voxelRenderer,
  clientId,
  getBlockType: () => selectedBlockType,
  getTickHint: () => snapshots.latest()?.tick ?? 0,
  sendCommand: (command) => snapshotSocket.sendCommand(command),
});
editController.connect();

window.addEventListener('beforeunload', () => {
  editController.disconnect();
  debugOverlay.dispose();
  entityRenderer.dispose();
  voxelRenderer.dispose();
  snapshotSocket.close();
});

window.wirecraft = {
  placeBlock(position: Position, blockType: BlockType = BlockType.DebugMover) {
    snapshotSocket.sendCommand({
      type: 'place_block',
      clientId,
      commandId: crypto.randomUUID(),
      tickHint: snapshots.latest()?.tick ?? 0,
      position,
      blockType,
    });
  },
  removeBlock(position: Position) {
    snapshotSocket.sendCommand({
      type: 'remove_block',
      clientId,
      commandId: crypto.randomUUID(),
      tickHint: snapshots.latest()?.tick ?? 0,
      position,
      blockType: BlockType.Air,
    });
  },
  snapshots,
};

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

function animate(timestampMs: number) {
  latestFps = calculateFps(previousFrameTimestampMs, timestampMs) ?? latestFps;
  previousFrameTimestampMs = timestampMs;
  scene.rotation.y += (targetYaw - scene.rotation.y) * 0.035;

  const renderServerTimeMs = snapshots.renderServerTimeMs(Date.now(), DEFAULT_INTERPOLATION_DELAY_MS);
  if (renderServerTimeMs !== null) {
    entityRenderer.updateFromSnapshots(snapshots.all(), renderServerTimeMs);
    app.dataset.renderServerTimeMs = String(Math.round(renderServerTimeMs));
    app.dataset.renderedEntities = String(entityRenderer.count);
  }

  debugOverlay.update({
    wsStatus,
    serverTick,
    snapshotBuffer: snapshots.length,
    renderedEntities: entityRenderer.count,
    fps: latestFps,
  });
  app.dataset.fps = latestFps === null ? '' : String(latestFps);
  app.dataset.debugVoxelBlocks = String(voxelBlocks);

  renderer.render(scene, camera);
  requestAnimationFrame(animate);
}

requestAnimationFrame(animate);
