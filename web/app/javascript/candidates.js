const API = "http://localhost:8080"

const list = document.getElementById("list")

async function loadCandidates() {

  const res = await fetch(`${API}/candidates`)
  const data = await res.json()

  list.innerHTML = ""

  data.forEach(c => {

    const div = document.createElement("div")

    div.className = "bg-slate-900 p-4 rounded flex justify-between items-center"

    div.innerHTML = `
      <div>
        <div class="font-semibold">${c.name}</div>
        <div class="text-xs text-slate-400">${c.email}</div>
      </div>

      <div class="flex gap-2">

        <button onclick="generateToken(${c.id})"
          class="bg-green-600 px-3 py-1 rounded text-sm">
          Get Link
        </button>

      </div>
    `

    list.appendChild(div)
  })
}

document.getElementById("createBtn").onclick = async () => {

  const name = document.getElementById("name").value
  const email = document.getElementById("email").value

  await fetch(`${API}/candidates`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ name, email })
  })

  loadCandidates()
}

async function generateToken(id) {

  const res = await fetch(`${API}/tokens`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      candidate_id: id
    })
  })

  const data = await res.json()

  const url = `${window.location.origin}/app/quiz.html?token=${data.token}`

  navigator.clipboard.writeText(url)

  alert("Link copied 🚀")
}

loadCandidates()
