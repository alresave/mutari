const API = "http://localhost:8080"

const token = new URLSearchParams(window.location.search).get("token")

const container = document.getElementById("questions")
const progress = document.getElementById("progress")

let answers = {}
let totalQuestions = 0

async function loadQuiz() {

  const res = await fetch(`${API}/quiz?token=${token}`)
  const data = await res.json()

  totalQuestions = data.questions.length

  renderQuestions(data.questions)
}

function renderQuestions(questions) {

  container.innerHTML = ""

  let currentBlock = null

  questions.forEach((q, index) => {

    // 🔥 block header
    if (q.block_name !== currentBlock) {

      currentBlock = q.block_name

      const header = document.createElement("div")
      header.className = "mt-8 mb-4 text-sm text-indigo-400 uppercase"

      header.innerText = currentBlock

      container.appendChild(header)
    }

    const div = document.createElement("div")
    div.className = "mb-6"

    div.innerHTML = `
      <div class="mb-3 font-medium text-white">
        ${index + 1}. ${q.question}
      </div>

      <div class="flex gap-2">
        ${[1,2,3,4,5].map(v => `
          <button
            data-q="${q.id}"
            data-v="${v}"
            class="option px-3 py-1 bg-slate-800 rounded hover:bg-indigo-600"
          >
            ${v}
          </button>
        `).join("")}
      </div>
    `

    container.appendChild(div)
  })

  bindEvents()
}

function bindEvents() {

  document.querySelectorAll(".option").forEach(btn => {

    btn.onclick = () => {

      const q = btn.dataset.q
      const v = parseInt(btn.dataset.v)

      answers[q] = v

      // 🔥 reset estilos
      document.querySelectorAll(`[data-q="${q}"]`).forEach(b => {
        b.classList.remove("bg-indigo-600")
        b.classList.add("bg-slate-800")
      })

      // 🔥 marcar seleccionado
      btn.classList.remove("bg-slate-800")
      btn.classList.add("bg-indigo-600")

      updateProgress()
    }
  })
}

function updateProgress() {

  const answered = Object.keys(answers).length

  progress.innerText = `${answered} / ${totalQuestions}`
}

document.getElementById("submitBtn").onclick = async () => {

  if (Object.keys(answers).length !== totalQuestions) {
    alert("Please answer all questions")
    return
  }

  const payload = Object.entries(answers).map(([q, v]) => ({
    question_id: parseInt(q),
    value: v
  }))

  const res = await fetch(`${API}/mcp/answers`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      token,
      answers: payload
    })
  })

  if(res.ok){
    document.body.innerHTML = `
      <div class="h-screen flex items-center justify-center text-center">
        <div>
          <div class="text-2xl font-bold mb-2">✅ Done</div>
          <div class="text-slate-400">Your assessment has been submitted</div>
        </div>
      </div>
    `
  } else {
    alert("Error submitting")
  }
}

loadQuiz()
