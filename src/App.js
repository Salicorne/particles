import { useState } from 'react';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Modal from 'react-bootstrap/Modal';
import Stack from 'react-bootstrap/Stack';
import InputGroup from 'react-bootstrap/InputGroup';
import Form from 'react-bootstrap/Form';
import Tabs from 'react-bootstrap/Tabs';
import Tab from 'react-bootstrap/tab';

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
        <h3>How does it work ?</h3>
        <p>The frontend interface is written in Javascript with <a href="https://react.dev/" target="_blank" rel="noreferrer">React</a>, mainly to manage the settings of the simulation. It instanciates and communicates with a Web Assembly application, written in Go. This application manages the particles simulation, and draws them on a Canvas on the page. </p>
      </Modal.Body>
      <Modal.Footer className="justify-content-between">
        <Stack direction='horizontal' gap={1}>
          <Button href="https://github.com/Salicorne/particles"><span className="bi bi-github"></span></Button>
          <Button href="https://vimeo.com/222974687" target="_blank"><span className="bi bi-vimeo"></span></Button>
        </Stack>
        <Button variant="primary" onClick={handleClose}>Ok</Button>
      </Modal.Footer>
    </Modal>
  );
}

function randomColor() {
  return "#"+(Math.floor(Math.random()*16777215)).toString(16);
}

//todo: avoid ascii attacks, use global counter in base52 with several letters instead
function nextPopulationId(populations) {
  let maxChar = "a".charCodeAt(0);
  for(let i = 0; i < Object.keys(populations).length; i++) {
    if(Object.keys(populations)[i].charCodeAt(0) > maxChar) {
      maxChar = Object.keys(populations)[i].charCodeAt(0);
    }
  }
  return String.fromCharCode([maxChar+1]);
}

function PopulationList({ populations, setPopulations }) {
/*
  function onPopulationColorChange(idx, newColor) {
    setPopulations(Object.keys(populations).map((populationName, i) => {
      if (i === idx) {
        return { n: populations[populationName].n, c: newColor };
      }
      return populations[populationName];
    }))
  }
*/
  function onPopulationColorChange(populationName, newColor) {
    let newPopulations = {};
    for(let [k, v] of Object.entries(populations)) {
      if (k === populationName) {
        newPopulations[k] = { n: populations[populationName].n, c: newColor };
      } else {
        newPopulations[k] = v;
      }
    }
    setPopulations(newPopulations);
  }
  function onPopulationNumberChange(populationName, newNumber) {
    let newPopulations = {};
    for(let [k, v] of Object.entries(populations)) {
      if (k === populationName) {
        newPopulations[k] = { n: newNumber, c: populations[populationName].c };
      } else {
        newPopulations[k] = v;
      }
    }
    setPopulations(newPopulations);
  }
  function deletePopulation(populationName) {
    let newPopulations = {};
    for(let i of Object.keys(populations)) {
      if (i !== populationName) {
        newPopulations[i] = populations[i];
      } 
    }
    setPopulations(newPopulations);
  }

  function addPopulation() {
    let newPopulations = {...populations};
    newPopulations[nextPopulationId(populations)] = {n: 10, c: randomColor()};
    setPopulations(newPopulations)
  }

  return (
    <Form>
      {Object.keys(populations).map((populationName, i) =>
        <Form.Group key={populationName} className="mb-3" >
          <InputGroup>
            <Form.Control column="lg" type="color" value={populations[populationName].c} onChange={(e) => onPopulationColorChange(populationName, e.target.value)} title="Choose the population color" />
            <Form.Control column="lg" type='number' value={populations[populationName].n} onChange={(e) => onPopulationNumberChange(populationName, parseInt(e.target.value))} />
            <Button variant="danger" onClick={(e) => deletePopulation(populationName)}>
              <span className="bi bi-trash" />
            </Button>
          </InputGroup>
        </Form.Group>
      )}
      <Button variant="primary" onClick={addPopulation} disabled={Object.keys(populations).length > 5}>
        +
      </Button>
    </Form>
  );
}

function SettingsModal({ show, handleClose, settings, setAppSettingsStr }) {
  let appSettings = settings !== "" ? JSON.parse(settings) : { p: {}, r: [] };
  let [populations, setPopulations] = useState(appSettings.p);

  function saveSettings() {
    setAppSettingsStr(JSON.stringify({ p: populations, r: appSettings.r }));
    window.setSettings(JSON.stringify({ p: populations, r: appSettings.r }));
    handleClose();
  }

  return (
    <Modal show={show} onHide={handleClose}>
      <Modal.Header closeButton>
        <Modal.Title>Settings</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Tabs defaultActiveKey="populations" id="uncontrolled-tab-example" className="mb-3">
          <Tab eventKey="populations" title="Populations">
          <PopulationList populations={populations} setPopulations={setPopulations} />
          </Tab>
          <Tab eventKey="rules" title="Rules">
            todo
          </Tab>
        </Tabs>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={handleClose}>Cancel</Button>
        <Button variant="primary" onClick={saveSettings}>Save</Button>
      </Modal.Footer>
    </Modal>
  )
}

let wasmInitialized = false;

function App() {
  //* Initialize React state
  const [showSettings, setShowSettings] = useState(false);
  const handleCloseSettings = () => setShowSettings(false);
  const handleShowSettings = () => setShowSettings(true);

  const [showInfo, setShowInfo] = useState(false);
  const handleCloseInfo = () => setShowInfo(false);
  const handleShowInfo = () => setShowInfo(true);

  const [appSettingsStr, setAppSettingsStr] = useState("")

  //* Initialize WASM app
  if (!wasmInitialized) {
    require("./wasm_exec");
    const go = new window.Go();
    WebAssembly.instantiateStreaming(fetch(process.env.PUBLIC_URL + "/app.wasm", { cache: "no-store" }), go.importObject).then((result) => {
      go.run(result.instance);
      window.startApp();
      setAppSettingsStr(window.getSettings())
      wasmInitialized = true;
    })
  }

  if (appSettingsStr === "") {
    return <p>Loading...</p>
  }

  return (
    <>
      <Container id="ui">

        <Stack direction="horizontal" gap={3}>
          <Button className='bi bi-gear-fill' onClick={handleShowSettings}></Button>
          <Button className='bi bi-info' onClick={handleShowInfo}></Button>
        </Stack>

        <pre className='font-monospace'>debug: {JSON.stringify(JSON.parse(appSettingsStr), null, 2)}</pre>

        {/* Settings modal */}
        <SettingsModal show={showSettings} handleClose={handleCloseSettings} settings={appSettingsStr} setAppSettingsStr={setAppSettingsStr} />

        {/* Info modal */}
        <InfoModal show={showInfo} handleClose={handleCloseInfo} />

      </Container>
    </>
  );
}

export default App;
