import {useState, useEffect, useRef} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";

// Terminal
import {Terminal} from "@xterm/xterm";
import { AttachAddon } from '@xterm/addon-attach';
import '@xterm/xterm/css/xterm.css';

function App() {
    //-------------------
    // Variables
    //-------------------

    // Template
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    // Terminal
    const term = new Terminal();
    const ref_term = useRef<HTMLDivElement>(null);
    const socket = new WebSocket("ws://localhost:8080");
    const attachAddon = new AttachAddon(socket);
    var term_loaded = false;

    //-------------------
    // Functions
    //-------------------
    function greet() {
        Greet(name).then(updateResultText);
    }

    function Block()
    {
        return (
            <div><h1>Test</h1></div>
        );
    }

    function WidgetList()
    {
        const data = ["hello", "bla", "test"];
        const listItems = data.map(t => <li>{t}</li>);
        const list = <ul>{listItems}</ul>;
        return list;
    }

    useEffect(() => {
        if (!term_loaded && ref_term.current)
        {
            term.open(ref_term.current);
            term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
            term.loadAddon(attachAddon);

            term_loaded = true;
        }
    });

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
            </div>
            <Block />
            <Block />
            <WidgetList />
            <div ref={ref_term}></div>
        </div>
    )
}

export default App
