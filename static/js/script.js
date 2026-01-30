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
  closeSidebar();
  CURRENT = "";
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

function sendMessage() {
  let jwt = localStorage.getItem("JWT");
  let msg = $("#msg").val();
  if (msg != "" && CURRENT != "") {
    let datos = {
      group: CURRENT,
      msg: msg,
    };
    jsonRequestNoDim(
      JSON.stringify(datos),
      "/api/sendMessage",
      jwt,
      (data, err) => {
        if (err != null) {
          showNotification(data.message, "danger");
          return;
        }
        $("#msg").val("");
        loadGroupChat2();
        scrollToBottom(1000);
        $("#msg").focus();
      },
    );
  }
}

var CURRENT = "";
var USERNAME = "";
var CONT = 0;

function getGroupsList() {
  let jwt = localStorage.getItem("JWT");
  jsonRequestNoDim(null, "/api/getGroupsList", jwt, (data, err) => {
    if (err != null) {
      showNotification(data.message, "danger");
      return;
    }
    let text = `
      <li class="group-item" onclick="createGroup()">
          <div class="group-avatar">++</div>
          <div class="group-info">
            <span class="group-name">CREAR NUEVO GRUPO</span>
          </div>
        </li>
      `;
    USERNAME = data.me;
    for (let i in data.data) {
      let unread = "";
      console.log(data.data[i].members[USERNAME]);
      if (data.data[i].members[USERNAME] == "UNREAD") {
        unread = `<span class="badge badge-warning">Msgs</span>`;
      }
      text += `
        <li class="group-item" onclick="loadGroupChat('${i}')">
          <div class="group-avatar">GRP</div>
          <div class="group-info">
            <span class="group-name">${data.data[i].name}</span>
            ${unread}
          </div>
        </li>
        `;
    }
    $("#groups").html(text);
  });
}

const BASEURL = "http://127.0.0.1:5000/app/";
var WAIT = 0;

function loadGroupChat(e) {
  closeSidebar();
  let jwt = localStorage.getItem("JWT");
  CURRENT = e;
  let datos = {
    msgs: e,
    max: 300,
  };
  WAIT = 1;
  $("#ginfo").html(`
         <h2>${e}</h2>
          <p>⏳ miembros</p>
        `);
  $("#menuBar").html("⏳ Info...");
  $("#chatMessages").html(
    `<span class="badge" style="font-size: 200px;">⏳</span>`,
  );
  jsonRequestNoDim(
    JSON.stringify(datos),
    "/api/loadGroupChat",
    jwt,
    (data, err) => {
      if (err != null) {
        showNotification(data.message, "danger");
        return;
      }
      let text = "";
      $("#chatMessages").html("");
      USERNAME = data.me;
      //console.log(data);
      for (let i in data.data) {
        //console.log(data.data[i]);
        let unread = "✔";
        if (data.data[i].from_user == data.me) {
          if (data.data[i].resents == 1) {
            unread = "✔✔";
          }
          let editmsgp = `<a href="javascript:void(0)" onclick="editMsg('${data.data[i].id}')">${data.data[i].text}</a>`;
          text += `
           <div class="message message-sent">
            ${editmsgp}
            <div class="message-status">${unread}</div>
            </div>
        `;
        } else {
          let editMember = data.data[i].apodo;
          let editmsg = data.data[i].text;
          if (data.datosGrupo.owner[data.me] !== undefined) {
            editMember = `<a href="javascript:void(0)" onclick="editMember('${data.data[i].from_user}', '${editMember}')">${data.data[i].apodo}</a>`;
            editmsg = `<a href="javascript:void(0)" onclick="editMsg('${data.data[i].id}')">${data.data[i].text}</a>`;
          }
          text += `
         <div class="message message-received">
            <div class="message-username">${editMember}</div>
            ${editmsg}
          </div>
        `;
        }
      }
      let lockd = `<button class="btn btn-warning" onclick="cerrarGrupo('${data.datosGrupo.id}')">Cerrar</button>`;
      if (data.datosGrupo.is_blocked == true) {
        lockd = `<button class="btn btn-primary" onclick="cerrarGrupo('${data.datosGrupo.id}')">Abrir</button>`;
      }
      let ttx = `
      <a href="javascript:void(0)" onclick="copy('${BASEURL + data.datosGrupo.link}')"><small>Enlace</small></a>
      ${lockd}
      <button class="btn btn-danger">ELIMINAR</button>
      `;
      if (data.datosGrupo.link == "") {
        ttx = `<button class="btn btn-danger" onclick="salirGrupo('${data.datosGrupo.id}')">SALIR DEL GRUPO</button>`;
      }

      if (
        data.datosGrupo.locked == true ||
        data.datosGrupo.is_blocked == true ||
        data.datosGrupo.members[USERNAME] == "SUSPENDED"
      ) {
        $("#msg").prop("disabled", true);
      } else {
        $("#msg").prop("disabled", false);
      }

      $("#chatMessages").html(text);
      $("#menuBar").html(ttx);
      $("#ginfo").html(`
         <h2>${data.datosGrupo.name}</h2>
          <p>${Object.keys(data.datosGrupo.members).length} miembros</p>
        `);
      scrollToBottom(1000);
      WAIT = 0;
    },
  );
}

function editMember(e, nombre) {
  $("#datos").html(`
    <p>Acciones sobre: ${nombre}</p>
    <button type="button" class="btn btn-danger" onclick="expulsar('${e}')">Expulsar</button> <button type="button" class="btn" onclick="suspender('${e}')"><span class="badge badge-warning">Suspender</span> | <span class="badge badge-success">Activar</span></button>`);
  openModal();
}

function editMsg(idmsg) {
  $("#datos").html(`
    <button type="button" class="btn btn-danger" onclick="eliminarMsg('${idmsg}')">ELIMINAR</button>`);
  openModal();
}

function salirGrupo(e) {
  let jwt = localStorage.getItem("JWT");
  let datos = {
    group: CURRENT,
  };
  jsonRequest(JSON.stringify(datos), "/api/salirGrupo", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      showNotification("error", "danger");
      return;
    }
    showNotification("ok", "success");
    CURRENT = "";
    $("#chatMessages").html("");
    $("#menuBar").html("");
    $("#ginfo").html("");
  });
}

function eliminarMsg(e) {
  let jwt = localStorage.getItem("JWT");
  let datos = {
    msg: e,
    group: CURRENT,
  };
  jsonRequest(JSON.stringify(datos), "/api/eliminarMsg", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      showNotification(err.message, "danger");
      return;
    }
    showNotification("Operacion exitosa", "success");
    closeModal();
    loadGroupChat2();
  });
}

function suspender(e) {
  let jwt = localStorage.getItem("JWT");
  let datos = {
    user: e,
    group: CURRENT,
  };
  jsonRequest(JSON.stringify(datos), "/api/suspender", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      showNotification(err.message, "danger");
      return;
    }
    showNotification("Operacion exitosa", "success");
    closeModal();
    //loadGroupChat2();
  });
}

function expulsar(e) {
  let jwt = localStorage.getItem("JWT");
  let datos = {
    user: e,
    group: CURRENT,
  };
  jsonRequest(JSON.stringify(datos), "/api/expulsar", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      showNotification(err.message, "danger");
      return;
    }
    showNotification("Operacion exitosa", "success");
    closeModal();
    //loadGroupChat2();
  });
}

var CHATS = {};

function loadGroupChat2() {
  if (CURRENT == "") {
    $("#ginfo").html("");
    return;
  }
  if (WAIT == 1) {
    return;
  }
  let jwt = localStorage.getItem("JWT");
  let datos = {
    msgs: CURRENT,
    max: 300,
  };
  jsonRequestNoDim(
    JSON.stringify(datos),
    "/api/loadGroupChat",
    jwt,
    (data, err) => {
      if (err != null) {
        showNotification(data.message, "danger");
        return;
      }
      let text = "";
      //$("#chatMessages").html();
      USERNAME = data.me;
      //console.log(data);
      let ctn = 0;
      for (let i in data.data) {
        ctn++;
        //console.log(data.data[i]);
        let unread = "✔";
        if (data.data[i].from_user == data.me) {
          if (data.data[i].resents == 1) {
            unread = "✔✔";
          }
          let editmsgp = `<a href="javascript:void(0)" onclick="editMsg('${data.data[i].id}')">${data.data[i].text}</a>`;
          text += `
           <div class="message message-sent">
            ${editmsgp}
            <div class="message-status">${unread}</div>
            </div>
        `;
        } else {
          let editMember = data.data[i].apodo;
          let editmsg = data.data[i].text;
          if (data.datosGrupo.owner[data.me] !== undefined) {
            editMember = `<a href="javascript:void(0)" onclick="editMember('${data.data[i].from_user}', '${editMember}')">${data.data[i].apodo}</a>`;
            editmsg = `<a href="javascript:void(0)" onclick="editMsg('${data.data[i].id}')">${data.data[i].text}</a>`;
          }
          text += `
         <div class="message message-received">
            <div class="message-username">${editMember}</div>
            ${editmsg}
          </div>
        `;
        }
      }

      let lockd = `<button class="btn btn-warning" onclick="cerrarGrupo('${data.datosGrupo.id}')">Cerrar</button>`;
      if (data.datosGrupo.is_blocked == true) {
        lockd = `<button class="btn btn-primary" onclick="cerrarGrupo('${data.datosGrupo.id}')">Abrir</button>`;
      }
      let ttx = `
      <a href="javascript:void(0)" onclick="copy('${BASEURL + data.datosGrupo.link}')"><small>Enlace</small></a>
      ${lockd}
      <button class="btn btn-danger">ELIMINAR</button>
      `;
      if (data.datosGrupo.link == "") {
        ttx = `<button class="btn btn-danger" onclick="salirGrupo('${data.datosGrupo.id}')">SALIR DEL GRUPO</button>`;
      }

      if (
        data.datosGrupo.locked == true ||
        data.datosGrupo.is_blocked == true ||
        data.datosGrupo.members[USERNAME] == "SUSPENDED"
      ) {
        $("#msg").prop("disabled", true);
      } else {
        $("#msg").prop("disabled", false);
      }

      $("#menuBar").html(ttx);
      $("#ginfo").html(`
         <h2>${data.datosGrupo.name}</h2>
          <p>${Object.keys(data.datosGrupo.members).length} miembros</p>
        `);
      if (ctn != CONT) {
        $("#chatMessages").html(text);
        scrollToBottom(1000);
        CONT = ctn;
      }
      console.log(CONT, ctn);
    },
  );
}

function cerrarGrupo(e) {
  let jwt = localStorage.getItem("JWT");
  let datos = {
    group: e,
  };
  jsonRequest(JSON.stringify(datos), "/api/cerrarGrupo", jwt, (data, err) => {
    if (err != null) {
      console.log(err);
      showNotification(err.message, "danger");
      return;
    }
    showNotification("Operacion exitosa", "success");
    //loadGroupChat2();
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

setInterval(validateSession, 5000);
setInterval(loadGroupChat2, 1500);
setInterval(getGroupsList, 3000);
getGroupsList();
