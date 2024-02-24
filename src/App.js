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

function SettingsModal({ show, handleClose, settings }) {
  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Settings</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {settings}
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>Cancel</Button>
        <Button variant="primary" onClick={handleClose}>Save</Button>
      </Modal.Footer>
    </Modal>
  )
}

function App() {
  //* Initialize React state
  const [showSettings, setShowSettings] = useState(false);
  const handleCloseSettings = () => setShowSettings(false);
  const handleShowSettings = () => setShowSettings(true);

  const [showInfo, setShowInfo] = useState(false);
  const handleCloseInfo = () => setShowInfo(false);
  const handleShowInfo = () => setShowInfo(true);

  const [appSettings, setAppSettings] = useState("")

  //* Initialize WASM app
  require("./wasm_exec");
  const go = new window.Go();
  WebAssembly.instantiateStreaming(fetch("/app.wasm", { cache: "no-store" }), go.importObject).then((result) => {
    go.run(result.instance);
    window.startApp();
    setAppSettings(window.getSettings())
  })

  return (
    <>
      <Container id="ui">

        <Stack direction="horizontal" gap={3}>
          <Button className='bi bi-gear-fill' onClick={handleShowSettings}></Button>
          <Button className='bi bi-info' onClick={handleShowInfo}></Button>
        </Stack>

        <pre className='font-monospace'>debug: {appSettings === "" ? "{}" : JSON.stringify(JSON.parse(appSettings), null, 2)}</pre>

        {/* Settings modal */}
        <SettingsModal show={showSettings} handleClose={handleCloseSettings} settings={appSettings} />

        {/* Info modal */}
        <InfoModal show={showInfo} handleClose={handleCloseInfo} />

      </Container>
    </>
  );
}

export default App;
