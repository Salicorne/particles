import { useState } from 'react';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Stack from 'react-bootstrap/Stack';

import './App.css';

function InfoModal({ show, handleClose }) {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Info</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <h3>What is this ?</h3>
        <p>"Particles" is a Go implementation of the "Clusters" concept explained by <a href="https://vimeo.com/222974687" target="_blank" rel="noreferrer">Jeffrey Ventrella</a>). Basically, the user can create a set of populations defined by a color, each population has a quantity of individuals. Each individual has a 2d position in a constrained space. The user can then create rules affecting all individuals of a population &alpha;. A rule specifies that individuals of a population &beta; affects individuals of this population &alpha; (it can be the same) with a force <i>f</i>, and with an effect range. If <i>f</i> &gt; 0 then individuals of &beta; repel individuals of &alpha;, and attract them otherwise. It appears that even a few rules can create interesting results, and complex behaviours can emerge from this algorithm even though not strictly defined in the code. I found this idea inspiring, and this is my attempt to implement it !</p>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="primary" onClick={handleClose}>Ok</Button>
      </Modal.Footer>
    </Modal>
  )
}

function SettingsModal({ show, handleClose }) {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Settings</Modal.Title>
      </Modal.Header>
      <Modal.Body>

      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>Cancel</Button>
        <Button variant="primary" onClick={handleClose}>Save</Button>
      </Modal.Footer>
    </Modal>
  )
}

function App() {
  const [showSettings, setShowSettings] = useState(false);
  const handleCloseSettings = () => setShowSettings(false);
  const handleShowSettings = () => setShowSettings(true);

  const [showInfo, setShowInfo] = useState(false);
  const handleCloseInfo = () => setShowInfo(false);
  const handleShowInfo = () => setShowInfo(true);

  let appSettings = { "c": { "#FF0000": 20, "#0000FF": 30 }, "r": [{ "c1": "#0000FF", "c2": "#0000FF", "f": -0.1, "r": 400 }, { "c1": "#FF0000", "c2": "#FF0000", "f": 0.8, "r": 200 }, { "c1": "#FF0000", "c2": "#0000FF", "f": 0.5, "r": 400 }, { "c1": "#0000FF", "c2": "#FF0000", "f": -0.8, "r": 200 }] }

  const fetchWasm = () => {
    console.log("Fetxhing wasm app...")
  }

  return (
    <>
      <Container id="ui" data-bs-theme="dark">

        <Stack direction="horizontal" gap={3}>
          <Button className='bi bi-gear-fill' onClick={handleShowSettings}></Button>
          <Button className='bi bi-info' onClick={handleShowInfo}></Button>
        </Stack>

        <pre className='font-monospace'>debug: {JSON.stringify(appSettings, null, 2)}</pre>

        {/* Settings modal */}
        <SettingsModal show={showSettings} handleClose={handleCloseSettings} />

        {/* Info modal */}
        <InfoModal show={showInfo} handleClose={handleCloseInfo} />

      </Container>

      {/* Particles canvas */}
      <canvas id="canvas" onLoadedData={fetchWasm}>HTML5 canvas is not supported on this browser</canvas>
    </>
  );
}

export default App;
