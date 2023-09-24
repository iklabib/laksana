<script setup>
import { onMounted, ref, watch} from 'vue';
import { basicSetup, EditorView } from "codemirror"
import { keymap } from "@codemirror/view"
import { cpp } from "@codemirror/lang-cpp"
import { python } from "@codemirror/lang-python"
import { indentWithTab } from "@codemirror/commands"
import { loadPyodide }from "pyodide"
import pythonExample from '../examples/python?raw'
import cExample from '../examples/c?raw'
import zigExample from '../examples/zig?raw'

let view = null
let py = null;

onMounted(async () => {
    view = new EditorView({
    doc:'#include <stdio.h>\nint main(){\n printf("Hello World!\\n");\n return 0;\n}',
    extensions: [
      basicSetup, 
      cpp(), 
      python(), 
      keymap.of([indentWithTab]),
    ],
    parent: document.getElementById("editor")
  })

  py = await loadPyodide({
    indexURL: "../pyodide/",
    stdout: (s) => {codeRes.value += s},
  });
})


const IDLE = 0
const RUNNING = 1
const SUCCESS = 2
const ERROR = 3

const langs = ['C', 'Python', 'Rust', 'Zig']
let selectedLang = ref(langs[0])

const examples = {
  "C": cExample,
  "Python": pythonExample,
  "Zig": zigExample
}
let selectedExample = ref("C")

let runState = ref("RUN")
let runClass = ref("btn-info")
let codeRes = ref("")

let it = 0

const runStateHandle = () => {
  switch (it) {
    case IDLE:
      runState.value = "RUN"
      runClass.value = "btn-info"
      const archive = view.state.doc.toString()
      runCode(archive, selectedLang.value)
      break
    case RUNNING:
      runState.value = "RUNNING"
      runClass.value = "btn-warning"
      break
    case SUCCESS:
      runState.value = "SUCCESS"
      runClass.value = "btn-success"
      break
    default:
      runState.value = "ERROR"
      runClass.value = "btn-error"
      break
  }
}

let buildRes = null
let runRes = null
const runCode = (archive, type) => {
  codeRes.value = ""
  if (type == "Python") {
    py.runPython(archive)
    return
  }

  let formData = new FormData();
  formData.append("archive", archive)
  formData.append("type", type)
  fetch('http://localhost:5176/run', {
    method: "POST",
    body: formData,
  }).then(response => response.json().then(data => {
      console.log(data)
      runRes = data.Run
      buildRes = data.Build

      codeRes.value = data.Run.Stdout
  }))
  .catch(f => console.log(f))
}

watch(selectedExample, (curr, prev) => {
  if (curr != prev) {
    view.dispatch({
      changes: {from: 0, to: view.state.doc.length, insert: examples[curr]}
    });
  }
}) 

</script>

<template>
  <div id='container' class='container mx-auto py-2 h-full'>
    <div class="flex flex-row w-64">
      <div class="flex-1">
        <label class="label">
          <span class="label-text">Language</span>
        </label>
      </div>

      <div class="flex-1">
        <select v-model="selectedLang" class="select select-bordered select-sm">
          <option  v-for="lang in langs">{{lang}}</option>
        </select>
      </div>

      <div class="flex-1">
        <label class="label">
          <span class="label-text">Examples</span>
        </label>
      </div>

      <div class="flex-1">
        <select v-model="selectedExample" class="select select-bordered select-sm">
          <option  v-for="example in Object.keys(examples)">{{example}}</option>
        </select>
      </div>

      <div class="flex-1">
        <button @click="runStateHandle" 
                :class='["btn", "btn-sm", runClass]'>
          {{runState}}
        </button>
      </div>
    </div>
    <div id="editor" class='border-2 h-1/2'></div>
    <div class="form-control h-1/3">
      <label class="label">
          <span class="label-text">Output</span>
      </label>
      <textarea v-model="codeRes" class="textarea textarea-bordered w-full h-full">
      </textarea>
    </div>
    </div>
</template>

<style scoped>
</style>
