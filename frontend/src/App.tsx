import {useState, useEffect, useRef, Fragment} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";

// Terminal
import {Terminal} from "@xterm/xterm";
import { AttachAddon } from '@xterm/addon-attach';
import '@xterm/xterm/css/xterm.css';

type Error =
{
    error_message: string
};

type isPreorder<A> = {
    compare: (a: A, b: A) => Boolean
};

type StateFragment =
{
    name: string,
    param: string,
    readCmd: string

};

type Fragment<Param,Value> =
{
    name: string,
    readInstructions: Instruction<Param,Value>[],
    writeInstructions: Instruction<[Param,Value],{}>[],
};

type FragmentReference =
{
    fragmentName: string
}

type Instruction<Input,Output> =
{
    dependencies: FragmentInstance[],
    command: (input: Input) => string,
    parseValue: (a: string) => Output | Error,
    compareValue: (a: Output, b: Output) => boolean
};

type FragmentInstance =
{
    fragmentReference: FragmentReference,
    value: object
};
function fragmentInstance<Value>(name: string, value: Value) : FragmentInstance
{
    return {
        fragmentReference: { fragmentName: name },
        value: value as object
    };
}

// FragmentInstance state
type FragmentInstanceState =
{
    state?: string
    required?: string
}


/////////////////////////////////////////
// Enum definitions

// OS

type OS_Ubuntu =
{
    OS_name: "Ubuntu"
}

type OS_NixOS =
{
    OS_name: "NixOS"
}

type OS =
  | OS_Ubuntu
  | OS_NixOS;

const NixOS : OS = {OS_name: "NixOS"}
const Ubuntu : OS = {OS_name: "Ubuntu"}


// Firewall Verdict

type FirewallVerdict_Allow =
{
    FirewallVerdict_verdict: "Allow"
}

type FirewallVerdict_Deny =
{
    FirewallVerdict_verdict: "Deny"
}

type FirewallVerdict =
 | FirewallVerdict_Allow
 | FirewallVerdict_Deny

const Allow: FirewallVerdict = {FirewallVerdict_verdict: "Allow"}
const Deny: FirewallVerdict = {FirewallVerdict_verdict: "Deny"}


// End Enum definitions
/////////////////////////////////////////


/////////////////////////////////////////
// Struct definitions

type Port =
{
    port: number
}

type FirewallConfig =
{
    enabled: boolean,
    firewall_rules: [Port,FirewallVerdict][]
}



// End Struct definitions
/////////////////////////////////////////

function MyTest() {
    const fragment_OS: Fragment<{},OS> =
    {
        name: "OS",
        readInstructions: [
            {
                dependencies: [],
                command: _ => "uname -a",
                parseValue: (s: string) =>
                {
                    if (s.includes("NixOS"))
                    {
                        return NixOS;
                    }
                    else if (s.includes("Ubuntu"))
                    {
                        return Ubuntu;
                    }
                    else
                    {
                        return {error_message: "Unknown os: " + s}
                    }
                },
                compareValue: (a, b) => a.OS_name == b.OS_name
            }
        ],
        writeInstructions: [],
    };

    const fragment_Firewall : Fragment<{},FirewallConfig> = 
    {
        name: "Firewall",
        readInstructions:
        [
            {
                dependencies:
                [
                    fragmentInstance<OS>("OS", Ubuntu)
                ],
                command: _ => "ufw status",
                parseValue: _ => {return {error_message: "could not read"}},
                compareValue: (a, b) => false
            }
        ],
        writeInstructions:
        []
    };

    var configuration: FragmentInstance[] =
    [
        fragmentInstance<OS>("OS", Ubuntu),
        fragmentInstance<FirewallConfig>("Firewall", 
            {
                enabled: true,
                firewall_rules:
                [
                    [{port: 80}, Allow]
                ]
            }
        )
    ];

    // function resolveFragment(target: FragmentReference, fragments: Fragment<any,any>[]): Fragment<any,any>
    // {

    // }


    function isExecutable(i: Instruction<any,any>, state: FragmentInstance[]): Boolean
    {
        return false;
    }

    function execute(i: Instruction<any,any>)
    {

    }


    function getStatus(ref_target: FragmentReference, fragments: Fragment<any,any>[], state: FragmentInstance[]): String | Error
    {
        // const target = resolveFragment(ref_target, fragments);

        // const executableInstructions = target.readInstructions.filter(instr => isExecutable(instr, state));

        // if (executableInstructions.length == 0)
        // {
        //     return {error_message: "No executable instructions"};
        // }
        // else if (executableInstructions.length > 1)
        // {
        //     return {error_message: "There are multiple executable instructions"};
        // }
        // else
        // {

        // }

        return {error_message: "Not implemented"};
    }


    // const os_frag: Fragment =
    // {
    //     name: "OS",
    //     readInstruction: "uname -a",
    //     readParser: output => "linux",
    // };

    // const app_docker_frag: Fragment =
    // {
    //     name: "Docker",
    //     readCommand: "docker --version",
    //     readParser: output => output,
    // };

    // const frag_firewall

}

type SystemState = Map<FragmentReference,FragmentInstanceState>;

/////////////////////////////////////////
// UI
function SystemStateUI()
{
    const [systemState, setSystemState] = useState<SystemState>(new Map([]));


    const updateSystemState = (key: FragmentReference, value: FragmentInstanceState) => {
        setSystemState(map => new Map(map.set(key, value)));
    }

    function changeState()
    {
        setSystemState(map => new Map(
            map.set(
                {fragmentName: "OS"},
                {state: undefined, required: undefined}
            )
        ));
    }

    type CombinedFragment = 
    {
        fragment_ref: FragmentReference,
        state?: FragmentInstanceState
    }

    function FragmentInstanceUI(props: CombinedFragment)
    {
        const display_state = props.state?.state?.toString() ?? "Unknown";
        const display_required = props.state?.required?.toString() ?? "Unknown";
        return (
            <div>{props.fragment_ref.fragmentName}: is={display_state}, should={display_required}</div>
        );
    }


    const items = [];
    const list = <ul></ul>;
    return (
        <div>
            <h2>length: {systemState.keys.length}</h2>

            <ul>
            {[...systemState.keys()].map(k => (
            <li key={k.fragmentName}>
                <FragmentInstanceUI fragment_ref={k} state={systemState.get(k)}  ></FragmentInstanceUI>
                </li>
            ))}
            </ul>

            <button className="btn" onClick={changeState}>Greet</button>
            <button className="btn">Test</button>
        </div>
    );
}



/////////////////////////////////////////
// Main

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

    //-------------------
    // Functions
    //-------------------
    function greet() {
        Greet(name).then(updateResultText);
        term.write("ls\n");
    }

    function openTerm() {
        if (ref_term.current)
        {
            term.open(ref_term.current);
            term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $ ')
            term.loadAddon(attachAddon);
        }
    }

    function pasteText() {
        term.input("ls\n");
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

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet</button>
                <button className='btn' onClick={openTerm}>open term</button>
                <button className='btn' onClick={pasteText}>paste text</button>
            </div>
            <SystemStateUI />
            <Block />
            <Block />
            <WidgetList />
            <div ref={ref_term}></div>
        </div>
    )
}

export default App
