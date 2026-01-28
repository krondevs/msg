const openMenu = document.getElementById("openMenu");
const sidebar = document.getElementById("sidebar");
const overlay = document.getElementById("overlay");
const toggleSearch = document.getElementById("toggleSearch");
const searchBar = document.getElementById("searchBar");
const searchInput = document.getElementById("searchInput");
const searchMeta = document.getElementById("searchMeta");
const chatMessages = document.getElementById("chatMessages");
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");
const toggleMenu = document.getElementById("toggleMenu");
const menuBar = document.getElementById("menuBar");
const JWT = localStorage.getItem("JWT");

if (JWT == null) {
  location.replace("/login");
}

let chatMessagesHTML = $("#chatMessages").html();

function restoreChat() {
  $("#chatMessages").html(chatMessagesHTML);
  $(".chat-input").show();
}

function createGroup() {
  $(".chat-input").hide();
  $("#chatMessages").html(`
        <form id="forma4">
            <div class="create-group">
                <div class="cg-header">
                    <button class="cg-back" id="cgBack">←</button>
                    <div class="cg-title">
                        <h2>Crear grupo</h2>
                        <span>Agrega participantes</span>
                    </div>
                </div>
        
                <div class="cg-content">
                    <div class="cg-card">
                        <div class="cg-avatar">G<small>p</small></div>
                        <div style="flex:1">
                            <div class="cg-field">
                                <input class="cg-input" type="text" placeholder="Nombre del grupo" name="name">
                                <span class="cg-help">Máx. 25 caracteres</span>
                            </div>
                        </div>
                    </div>
        
                    <div class="cg-section">
                        <h3>Datos del grupo</h3>
        
                        <div class="cg-grid">
                            <div class="cg-field">
                                <label class="cg-label">Especialidad</label>
                                <select class="cg-input" name="specialty">
                                    <option>Arquitectura</option>
                                    <option>Ingeniería</option>
                                    <option>Diseño Interior</option>
                                    <option>Construcción</option>
                                    <option>Topografía</option>
                                </select>
                            </div>
                        </div>
        
                        <div class="cg-field">
                            <label class="cg-label">Descripción</label>
                            <textarea class="cg-input cg-textarea" rows="3" placeholder="Describe el objetivo del grupo"
                                name="description"></textarea>
                        </div>
        
                        <div class="cg-grid">
                            <div class="cg-field">
                                <label class="cg-label">País</label>
                                <input class="cg-input" type="text" placeholder="Ej. Colombia" name="country">
                            </div>
                            <div class="cg-field">
                                <label class="cg-label">Ciudad</label>
                                <input class="cg-input" type="text" placeholder="Ej. Bogotá" name="city">
                            </div>
                        </div>
        
                        <div class="cg-grid">
                            <div class="cg-field">
                                <label class="cg-label">Intermediario</label>
                                <input class="cg-input" type="text" placeholder="Nombre del intermediario" name="intermediary">
                            </div>
                            <div class="cg-field">
                                <label class="cg-label">Contacto</label>
                                <input class="cg-input" type="text" placeholder="Email o teléfono" name="contact">
                            </div>
                        </div>
                    </div>        
                    <div class="cg-actions">
                        <button class="cg-btn cg-ghost" id="cgCancel">Cancelar</button>
                        <button type="button" class="cg-btn cg-primary" id="cgCreate">Crear grupo</button>
                    </div>
                </div>
            </div>
        </form>
    `);

  $("#cgBack, #cgCancel").on("click", restoreChat);
  $("#cgCreate").on("click", addNewGroup);
}

function addNewGroup() {
  let datos = objetizeForm("forma4");
  datos.contact += "";
  let jwt = localStorage.getItem("JWT");
  jsonRequest(JSON.stringify(datos), "/api/addNewGroup", jwt, (data, err) => {
    if (err != null) {
      showNotification(err.message, "danger");
      return;
    }
    showNotification(data.message, "success");
    getGroupsList();
    closeModal();
    $("#cgBack").click();
  });
}

function getGroupsList() {
  let jwt = localStorage.getItem("JWT");
  jsonRequestNoDim(null, "/api/getGroupsList", jwt, (data, err) => {
    if (err != null) {
      showNotification(data.message, "danger");
      return
    }
    let text = $("#groups").html();
    for (let i in data.data) {
      text += `
            <li class="group-item">
            <div class="group-avatar">${data.data[i].avatar}</div>
            <div class="group-info">
                <span class="group-name">${data.data[i].name}</span>
                <span class="group-meta">${data.data[i].cant}</span>
            </div>
            </li>
           `;
      $("#groups").html(text);
    }
  });
}

function validateSession() {
  let jwt = localStorage.getItem("JWT");
  jsonRequestNoDim(null, "/api/validateSession", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      if (err.message == "token expired") {
        localStorage.clear();
        location.replace("/login");
      }
    }
  });
}

setInterval(validateSession, 5000);
getGroupsList();

let matchIndex = 0;
let matches = [];

function openSidebar() {
  sidebar.classList.add("open");
  overlay.classList.add("show");
}

function closeSidebar() {
  sidebar.classList.remove("open");
  overlay.classList.remove("show");
}

openMenu.addEventListener("click", openSidebar);
overlay.addEventListener("click", closeSidebar);
document.addEventListener("keydown", (e) => {
  if (e.key === "Escape") {
    closeSidebar();
    closeSearch();
  }
});

function closeSearch() {
  searchBar.classList.remove("show");
  chatMessages.classList.remove("with-search");
  searchInput.value = "";
  clearHighlights();
  updateMeta();
}

function closeSearch2() {
  menuBar.classList.remove("show");
}

toggleSearch.addEventListener("click", () => {
  const isOpen = searchBar.classList.contains("show");
  if (isOpen) {
    closeSearch();
  } else {
    searchBar.classList.add("show");
    chatMessages.classList.add("with-search");
    searchInput.focus();
  }
});

toggleMenu.addEventListener("click", () => {
  const isOpen = menuBar.classList.contains("show");
  if (isOpen) {
    closeSearch2();
  } else {
    menuBar.classList.add("show");
  }
});

function escapeRegExp(str) {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "$&");
}

function clearHighlights() {
  const highlighted = chatMessages.querySelectorAll("span.highlight");
  highlighted.forEach((span) => {
    const textNode = document.createTextNode(span.textContent);
    const parent = span.parentNode;
    parent.replaceChild(textNode, span);
    parent.normalize();
  });
  matches = [];
  matchIndex = 0;
}

function getTextNodes(root) {
  const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT, {
    acceptNode(node) {
      if (!node.nodeValue.trim()) {
        return NodeFilter.FILTER_REJECT;
      }
      const parent = node.parentElement;
      if (parent) {
        if (parent.classList.contains("message-username")) {
          return NodeFilter.FILTER_REJECT;
        }
        if (parent.classList.contains("message-status")) {
          return NodeFilter.FILTER_REJECT;
        }
      }
      return NodeFilter.FILTER_ACCEPT;
    },
  });

  const nodes = [];
  let current;
  while ((current = walker.nextNode())) {
    nodes.push(current);
  }
  return nodes;
}

function highlightMatches(query) {
  clearHighlights();
  if (!query) {
    updateMeta();
    return;
  }

  const q = query.toLowerCase();
  const messages = chatMessages.querySelectorAll(".message");

  messages.forEach((msg) => {
    const textNodes = getTextNodes(msg);

    textNodes.forEach((node) => {
      const text = node.nodeValue;
      const lower = text.toLowerCase();
      let index = 0;

      if (lower.indexOf(q) !== -1) {
        const frag = document.createDocumentFragment();

        while (true) {
          const found = lower.indexOf(q, index);
          if (found === -1) {
            break;
          }
          frag.append(document.createTextNode(text.slice(index, found)));
          const span = document.createElement("span");
          span.className = "highlight";
          span.textContent = text.slice(found, found + q.length);
          frag.append(span);
          index = found + q.length;
        }

        frag.append(document.createTextNode(text.slice(index)));
        node.replaceWith(frag);
      }
    });
  });

  matches = Array.from(chatMessages.querySelectorAll("span.highlight"));
  matchIndex = 0;
  updateMeta();
  scrollToMatch();
}

function updateMeta() {
  const hasMatches = matches.length > 0;
  searchMeta.textContent = hasMatches
    ? `${matchIndex + 1}/${matches.length}`
    : "0/0";
  prevBtn.disabled = !hasMatches;
  nextBtn.disabled = !hasMatches;
}

function scrollToMatch() {
  if (!matches.length) {
    return;
  }
  const target = matches[matchIndex];
  const parent = target.closest(".message");
  if (parent) {
    parent.scrollIntoView({ behavior: "smooth", block: "center" });
  }
}

function goNext() {
  if (!matches.length) {
    return;
  }
  matchIndex = (matchIndex + 1) % matches.length;
  updateMeta();
  scrollToMatch();
}

function goPrev() {
  if (!matches.length) {
    return;
  }
  matchIndex = (matchIndex - 1 + matches.length) % matches.length;
  updateMeta();
  scrollToMatch();
}

prevBtn.addEventListener("click", goPrev);
nextBtn.addEventListener("click", goNext);

searchInput.addEventListener("input", (e) => {
  highlightMatches(e.target.value.trim());
});

searchInput.addEventListener("keydown", (e) => {
  if (!matches.length) {
    return;
  }
  if (e.key === "Enter") {
    e.preventDefault();
    if (e.shiftKey) {
      goPrev();
    } else {
      goNext();
    }
  }
});

function scrollToBottom(duration = 300) {
  if (!chatMessages) {
    return;
  }

  const reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  if (reduce || duration <= 0) {
    chatMessages.scrollTop = chatMessages.scrollHeight;
    return;
  }

  const start = chatMessages.scrollTop;
  const end = chatMessages.scrollHeight;
  const change = end - start;
  const startTime = performance.now();

  function easeOutCubic(t) {
    return 1 - Math.pow(1 - t, 3);
  }

  function animate(now) {
    const elapsed = now - startTime;
    const progress = Math.min(elapsed / duration, 1);
    const eased = easeOutCubic(progress);
    chatMessages.scrollTop = start + change * eased;

    if (progress < 1) {
      requestAnimationFrame(animate);
    }
  }

  requestAnimationFrame(animate);
}

document.addEventListener("DOMContentLoaded", () => {
  scrollToBottom();
});
