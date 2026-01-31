/**
 * jqueryx - Librería de extensión para jQuery
 *
 * @library jqueryx
 * @version 1.0.0
 * @author Jober Urizare
 * @copyright Apix Technology
 *
 * @description
 * jqueryx es una librería de extensión para jQuery que proporciona
 * herramientas avanzadas para desarrollo web, inspirada en la filosofía
 * de HTMX, con funcionalidades para manejo de solicitudes AJAX,
 * notificaciones, modales, gráficos y utilidades de interfaz.
 *
 * @features
 * - Gestión de solicitudes AJAX simplificadas
 * - Componentes de interfaz dinámicos (modales, notificaciones)
 * - Generación de gráficos con Chart.js
 * - Generación de códigos QR
 * - Manejo de overlays y spinners de carga
 * - Serialización de formularios
 * - Gestión de colores dinámicos
 *
 * @dependencies
 * - jQuery (3.x o superior)
 * - Chart.js (para funciones de gráficos)
 * - QRCode.js (para generación de códigos QR)
 *
 * @modules
 * - Modal Management
 * - Notification System
 * - AJAX Utilities
 * - Loading Overlay
 * - Charting
 * - QR Code Generation
 *
 * @license MIT
 *
 * @contact
 * - Email: apixtechnologyca@gmail.com
 * - Web: https://apixtech.github.io/index
 *
 * @usage
 * // Ejemplos básicos de uso
 *
 * // Mostrar notificación
 * showNotification('Operación exitosa', 'success');
 *
 * // Abrir modal
 * openModal();
 *
 * // Enviar solicitud AJAX
 * jsonRequest(datos, '/ruta', (respuesta, error) => { ... });
 *
 * @roadmap
 * - Mejora de rendimiento
 * - Más componentes de interfaz
 * - Soporte para más tipos de gráficos
 *
 * @notes
 * - Diseñada para ser ligera y eficiente
 * - Fácil integración con proyectos existentes
 * - Inspirada en las mejores prácticas de desarrollo web
 */

function openModal() {
  const modal = document.getElementById("miModal");
  modal.style.display = "flex";

  // Pequeño retraso para permitir que el display: flex se aplique
  setTimeout(() => {
    modal.classList.add("mostrar");
  }, 10);
}

function uuid() {
  return typeof crypto !== "undefined" && crypto.randomUUID
    ? crypto.randomUUID()
    : ([1e7] + -1e3 + -4e3 + -8e3 + -1e11).replace(/[018]/g, (c) =>
        (
          c ^
          (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (c / 4)))
        ).toString(16),
      );
}

function closeModal() {
  const modal = document.getElementById("miModal");
  modal.classList.remove("mostrar");

  // Esperar a que termine la animación antes de ocultar
  setTimeout(() => {
    modal.style.display = "none";
  }, 300);
}

function cerrarModal() {
  const modal = document.getElementById("miModal");
  modal.classList.remove("mostrar");

  // Esperar a que termine la animación antes de ocultar
  setTimeout(() => {
    modal.style.display = "none";
  }, 300);
}

/*      // Cerrar modal si se hace clic fuera del contenido
window.onclick = function(event) {
	const modal = document.getElementById('miModal');
	 if (event.target == modal) {
		 closeModal();
	}
}*/

function serializeForm(formId) {
  const arr = $("#" + formId).serializeArray();
  console.log(arr);
  let obj = {};
  for (let item of arr) {
    obj[item.name] = item.value;
  }
  return JSON.stringify(obj);
}

function objetizeForm(formId) {
  const arr = $("#" + formId).serializeArray();
  console.log(arr);
  let obj = {};
  for (let item of arr) {
    item.value = item.value.trim();
    if (item.value !== "" && !isNaN(item.value)) {
      item.value = parseFloat(item.value);
    }
    obj[item.name] = item.value;
  }
  return obj;
}
/**
 * Realiza una solicitud AJAX con datos JSON y maneja la respuesta mediante un callback
 *
 * @param {string} jsonDataSerialized - Datos serializados en formato JSON
 * @param {string} route - URL del endpoint del servidor
 * @param {function} callback - Función de retorno con dos parámetros (data, error)
 *
 * @description
 * - Envía una solicitud POST con datos en formato JSON
 * - Establece el tipo de contenido como "application/json"
 * - Invoca el callback con los datos de respuesta o el error
 * - Registra errores en la consola para depuración
 * - Permite manejo flexible de respuestas y errores en el código del cliente
 *
 * @example
 * // Ejemplo de uso para enviar datos de usuario
 * const userData = JSON.stringify({
 *   nombre: 'Juan',
 *   email: 'juan@ejemplo.com'
 * });
 *
 * jsonRequest(
 *   userData,
 *   '/api/usuarios',
 *   (data, error) => {
 *     if (error) {
 *       console.error('Error:', error);
 *       return;
 *     }
 *     console.log('Respuesta:', data);
 *   }
 * )
 */
function jsonRequest(jsonDataSerialized, route, bearerToken, callback) {
  oscurecerPantalla();
  $.ajax({
    type: "POST",
    url: route,
    contentType: "application/json",
    data: jsonDataSerialized,
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
    success: (data) => {
      restablecerPantalla();
      callback(data, null);
    },
    error: function (xhr, status, error) {
      console.error(xhr, status, error);
      restablecerPantalla();
      let json;
      if (xhr.responseJSON) {
        json = xhr.responseJSON;
      } else {
        try {
          json = JSON.parse(xhr.responseText);
        } catch (_) {
          json = { status: "unknown", message: error };
        }
      }

      callback(null, json);
      //callback(null, xhr);
    },
  });
}

function jsonRequestNoDim(jsonDataSerialized, route, bearerToken, callback) {
  $.ajax({
    type: "POST",
    url: route,
    contentType: "application/json",
    data: jsonDataSerialized,
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
    success: (data) => {
      callback(data, null);
    },
    error: function (xhr, status, error) {
      console.error(xhr, status, error);
      let json;
      if (xhr.responseJSON) {
        json = xhr.responseJSON;
      } else {
        try {
          json = JSON.parse(xhr.responseText);
        } catch (_) {
          json = { status: "unknown", message: error };
        }
      }
      callback(null, json);
    },
  });
}

function sendMultiPartData(method, route, formId, callback) {
  oscurecerPantalla();
  let formData = new FormData(document.getElementById(formId));
  $.ajax({
    type: method,
    url: route,
    data: formData,
    processData: false,
    contentType: false,
    success: (data) => {
      restablecerPantalla();
      callback(data, null);
    },
    error: (xhr, status, error) => {
      console.error(xhr, status, error);
      restablecerPantalla();
      callback(null, error);
    },
  });
}

function sendMultiPartData2(method, route, bearerToken, formId, callback) {
  //oscurecerPantalla();
  let formData = new FormData(document.getElementById(formId));
  $.ajax({
    type: method,
    url: route,
    data: formData,
    processData: false,
    contentType: false,
    headers: {
      Authorization: `Bearer ${bearerToken}`,
    },
    success: (data) => {
      //restablecerPantalla();
      callback(data, null);
    },
    error: (xhr, status, error) => {
      console.error(xhr, status, error);
      //restablecerPantalla();
      callback(null, error);
    },
  });
}

function copy(text) {
  // Verificar si la API Clipboard está disponible (navegadores modernos)
  if (navigator.clipboard && window.isSecureContext) {
    // Usar la nueva API Clipboard
    navigator.clipboard
      .writeText(text)
      .then(function () {
        console.log("Texto copiado al portapapeles: " + text);
      })
      .catch(function (err) {
        console.error("Error al copiar: ", err);
      });
  } else {
    // Fallback para navegadores antiguos
    try {
      // Crear un elemento temporal de texto
      var textarea = document.createElement("textarea");
      textarea.value = text;
      textarea.style.position = "fixed";
      textarea.style.left = "-9999px";
      textarea.style.top = "-9999px";
      document.body.appendChild(textarea);

      // Seleccionar el texto
      textarea.select();
      textarea.setSelectionRange(0, text.length);

      // Intentar copiar
      var successful = document.execCommand("copy");

      // Eliminar el elemento temporal
      document.body.removeChild(textarea);

      if (successful) {
        console.log("Texto copiado al portapapeles: " + text);
      } else {
        console.error("Error al copiar el texto");
      }
    } catch (err) {
      console.error("Error al copiar: ", err);
    }
  }
  showNotification("Enlace copiado", "success");
}

function sendForm(method, route, formId, callback) {
  oscurecerPantalla();
  let formData = new FormData(document.getElementById(formId));
  $.ajax({
    type: method,
    url: route,
    data: formData,
    processData: false,
    contentType: false,
    success: (data) => {
      restablecerPantalla();
      callback(data, null);
    },
    error: (xhr, status, error) => {
      console.error(xhr, status, error);
      restablecerPantalla();
      callback(null, error);
    },
  });
}

function loadingIn() {
  // Crear un overlay si no existe
  if (!document.getElementById("loading-overlay")) {
    const overlay = document.createElement("div");
    overlay.id = "loading-overlay";
    overlay.style.position = "fixed";
    overlay.style.top = "0";
    overlay.style.left = "0";
    overlay.style.width = "100%";
    overlay.style.height = "100%";
    overlay.style.backgroundColor = "rgba(0, 0, 0, 0.5)"; // Semitransparente
    overlay.style.zIndex = "9999"; // Asegura que esté encima de otros elementos
    overlay.style.display = "flex";
    overlay.style.justifyContent = "center";
    overlay.style.alignItems = "center";

    // Opcional: Añadir un spinner de carga
    const spinner = document.createElement("div");
    spinner.style.width = "50px";
    spinner.style.height = "50px";
    spinner.style.border = "5px solid #f3f3f3";
    spinner.style.borderTop = "5px solid #3498db";
    spinner.style.borderRadius = "50%";
    spinner.style.animation = "spin 1s linear infinite";

    // Añadir estilo de animación
    const styleSheet = document.createElement("style");
    styleSheet.textContent = `
			@keyframes spin {
				0% { transform: rotate(0deg); }
				100% { transform: rotate(360deg); }
			}
		`;
    document.head.appendChild(styleSheet);

    overlay.appendChild(spinner);
    document.body.appendChild(overlay);
  }
}

function loadingOut() {
  const overlay = document.getElementById("loading-overlay");
  if (overlay) {
    overlay.remove();
  }
}

/**
 * Muestra una notificación emergente con estilo y tipo personalizado
 *
 * @param {string} message - Mensaje de texto a mostrar en la notificación
 * @param {string} [type='success'] - Tipo de notificación (success, warning, danger, info)
 * @param {number} [duration=5000] - Tiempo en milisegundos que se mostrará la notificación
 *
 * @description
 * - Crea dinámicamente un elemento de notificación
 * - Utiliza iconos predefinidos según el tipo de notificación
 * - Agrega la notificación a un contenedor específico
 * - Permite cierre manual o automático
 * - Soporta diferentes estilos (success, warning, danger, info)
 * - Añade transición de aparición suave
 *
 * @example
 * // Mostrar notificación de éxito
 * showNotification('Operación completada con éxito');
 *
 * // Mostrar notificación de advertencia con duración personalizada
 * showNotification('Advertencia: Datos no guardados', 'warning', 3000);
 *
 * // Mostrar notificación de error
 * showNotification('Error al procesar la solicitud', 'danger');
 */
function showNotification(message, type = "success", duration = 5000) {
  // Mapeo de iconos para cada tipo de notificación
  const icons = {
    success: "✔️",
    warning: "⚠️",
    danger: "❌",
    error: "❌",
    info: "ℹ️",
  };

  // Crear el elemento de notificación
  const notification = document.createElement("div");
  notification.classList.add("notification", type);

  // Contenido de la notificación
  notification.innerHTML = `
	  <div class="notification-icon">
		${icons[type] || "✔️"}
	  </div>
	  <div class="notification-content">
		${message}
	  </div>
	  <button class="notification-close" onclick="closeNotification(this)">
		&times;
	  </button>
	`;

  // Obtener el contenedor de notificaciones
  const container = document.getElementById("notificationContainer");

  // Agregar la notificación al contenedor
  container.appendChild(notification);

  // Mostrar la notificación
  setTimeout(() => {
    notification.classList.add("show");
  }, 10);

  // Ocultar y eliminar la notificación después de la duración especificada
  const timer = setTimeout(() => {
    closeNotification(notification);
  }, duration);

  // Añadir método para cerrar manualmente
  notification.closeTimer = timer;
}

/**
 * Cierra una notificación específica
 *
 * @param {HTMLElement} notificationElement - Elemento de notificación a cerrar
 */
function closeNotification(element) {
  // Obtener el elemento de notificación
  const notification = element.closest(".notification");

  // Limpiar cualquier temporizador existente
  if (notification.closeTimer) {
    clearTimeout(notification.closeTimer);
  }

  // Quitar la clase de mostrar
  notification.classList.remove("show");

  // Eliminar la notificación después de la animación
  setTimeout(() => {
    notification.remove();
  }, 300);
}

function makeCode(qrcod, text) {
  qrcod.makeCode(text);
}

function codeQrMaker(divId, text) {
  let qrcode = new QRCode(document.getElementById(divId), {
    width: 100,
    height: 100,
  });
  makeCode(qrcode, text);
  despues();
}

/**
 * Genera un color dinámico basado en el valor proporcionado y el valor máximo
 *
 * @param {number} valor - Valor numérico para generar el color
 * @param {number} maximo - Valor máximo para calcular la escala de colores
 * @returns {string} Color en formato rgba
 *
 * @description
 * - Transforma un valor numérico en un color que varía de rojo a verde
 * - Utiliza una interpolación lineal entre rojo (valores bajos) y verde (valores altos)
 * - Normaliza los colores basándose en el valor máximo proporcionado
 * - Devuelve un color con transparencia de 0.7
 *
 * @example
 * // Generar un color para un valor de 50 en un rango máximo de 100
 * const color = generarColor(50, 100);
 * // Resultado podría ser: 'rgba(127, 127, 0, 0.7)'
 *
 * @notes
 * - El valor se convierte a su valor absoluto para manejar números negativos
 * - El componente azul siempre es 0
 * - La transparencia es fija en 0.7
 */
function generarColor(valor, maximo) {
  // Calcular el valor de rojo, verde y azul en base al valor
  valor = Math.abs(valor);
  let kk;
  let cc;
  cc = 255 / maximo;
  let r = Math.floor(255 - valor * cc); // Valor más bajo, más cerca del rojo
  let g = Math.floor(valor * cc); // Valor más alto, más cerca del verde
  let b = 0; // Sin componente azul
  return "rgba(" + r + ", " + g + ", " + b + ", 0.7)";
}

/**
 * Crea un gráfico de líneas utilizando Chart.js
 *
 * @param {string[]} leyendas - Etiquetas para el eje X del gráfico
 * @param {number[]} valores - Valores de datos correspondientes a las etiquetas
 * @param {string} titulo - Título del gráfico
 * @param {string} canvasId - ID del elemento canvas donde se renderizará el gráfico
 *
 * @description
 * - Genera un gráfico de líneas con un conjunto de datos
 * - Configura estilos personalizados para ejes X e Y
 * - Permite personalización de colores, tamaños de fuente y título
 * - Utiliza Chart.js para la renderización del gráfico
 *
 * @example
 * // Crear un gráfico de ventas mensuales
 * const meses = ['Enero', 'Febrero', 'Marzo', 'Abril'];
 * const ventas = [1000, 1500, 1200, 1800];
 * chartLine(meses, ventas, 'Ventas Mensuales 2023', 'graficoVentas');
 *
 * @requires Chart.js - Librería de visualización de gráficos
 */
function chartLine(leyendas, valores, titulo, canvasId) {
  let $grafica = document.querySelector("#" + canvasId);
  // Las etiquetas son las que van en el eje X.
  let etiquetas = leyendas;
  // Podemos tener varios conjuntos de datos. Comencemos con uno
  let datosVentas2020 = {
    label: titulo,
    data: valores, // La data es un arreglo que debe tener la misma cantidad de valores que la cantidad de etiquetas
    backgroundColor: "rgba(54, 162, 235, 0.2)", // Color de fondo
    borderColor: "#000", // Color del borde
    borderWidth: 2, // Ancho del borde
  };
  new Chart($grafica, {
    type: "line", // Tipo de gráfica
    data: {
      labels: etiquetas,
      datasets: [
        datosVentas2020,
        // Aquí más datos...
      ],
    },
    options: {
      scales: {
        yAxes: [
          {
            ticks: {
              beginAtZero: true,
              fontSize: 20, // Tamaño de letra en el eje X
              fontColor: "black",
            },
          },
        ],
        xAxes: [
          {
            ticks: {
              fontSize: 20, // Tamaño de letra en el eje X
              fontColor: "black",
            },
          },
        ],
      },
      title: {
        display: true,
        text: titulo,
        fontSize: 24, // Tamaño de letra del título de la gráfica
        fontColor: "black",
      },
    },
  });
}

/**
 * Crea un gráfico de barras utilizando Chart.js con colores dinámicos
 *
 * @param {string[]} leyendas - Etiquetas para el eje X del gráfico
 * @param {number[]} valores - Valores numéricos para cada barra
 * @param {string} titulo - Título del gráfico
 * @param {string} canvasId - ID del elemento canvas donde se renderizará el gráfico
 *
 * @description
 * - Genera un gráfico de barras con colores dinámicos basados en los valores
 * - Calcula el valor máximo para generar una escala de colores
 * - Configura estilos personalizados para ejes X e Y
 * - Utiliza Chart.js para la renderización del gráfico
 *
 * @example
 * // Crear un gráfico de barras de ventas por categoría
 * const categorias = ['Electrónica', 'Ropa', 'Alimentos', 'Deportes'];
 * const ventas = [5000, 3500, 2800, 4200];
 * chartBar(categorias, ventas, 'Ventas por Categoría', 'graficoVentas');
 *
 * @requires Chart.js - Librería de visualización de gráficos
 * @requires generarColor() - Función personalizada para generar colores dinámicos
 */
function chartBar(leyendas, valores, titulo, canvasId) {
  // leyendas son slices []
  //valores son slices[]

  let $grafica = document.querySelector("#" + canvasId);
  let etiquetas = leyendas;
  let maximo = Math.max(...valores);
  let datosVentas2020 = {
    label: titulo,
    data: valores,
    backgroundColor: valores.map((valor) => generarColor(valor, maximo)), // Generar colores dinámicamente
    borderColor: "rgba(0, 0, 0, 1)", // Color del borde de las barras
    borderWidth: 2,
  };

  new Chart($grafica, {
    type: "bar",
    data: {
      labels: etiquetas,
      datasets: [datosVentas2020],
    },
    options: {
      scales: {
        yAxes: [
          {
            ticks: {
              beginAtZero: true,
              fontSize: 20,
              fontColor: "black",
            },
          },
        ],
        xAxes: [
          {
            ticks: {
              fontSize: 20,
              fontColor: "black",
            },
          },
        ],
      },
      title: {
        display: true,
        text: titulo,
        fontSize: 24,
        fontColor: "black",
      },
    },
  });
}

/**
 * Oscurece la pantalla y muestra un indicador de carga
 *
 * @function oscurecerPantalla
 * @description
 * - Crea una capa semitransparente que cubre toda la pantalla
 * - Añade un spinner de carga centrado en la pantalla
 * - Previene la creación múltiple de overlays
 *
 * @returns {void}
 *
 * @behavior
 * - Si ya existe un overlay, no crea uno nuevo
 * - Bloquea interacciones con el contenido subyacente
 * - Añade un spinner animado giratorio
 *
 * @example
 * // Mostrar overlay de carga
 * oscurecerPantalla();
 *
 * @notes
 * - Utiliza z-index alto para asegurar visibilidad
 * - Centrado tanto horizontal como verticalmente
 * - Añade animación de rotación al spinner
 *
 * @warning
 * - Solo debe llamarse una vez
 * - Para remover el overlay, será necesaria una función complementaria
 */
function oscurecerPantalla() {
  // Crear un overlay si no existe
  if (!document.getElementById("loading-overlay")) {
    const overlay = document.createElement("div");
    overlay.id = "loading-overlay";
    overlay.style.position = "fixed";
    overlay.style.top = "0";
    overlay.style.left = "0";
    overlay.style.width = "100%";
    overlay.style.height = "100%";
    overlay.style.backgroundColor = "rgba(0, 0, 0, 0.5)"; // Semitransparente
    overlay.style.zIndex = "9999"; // Asegura que esté encima de otros elementos
    overlay.style.display = "flex";
    overlay.style.justifyContent = "center";
    overlay.style.alignItems = "center";

    // Opcional: Añadir un spinner de carga
    const spinner = document.createElement("div");
    spinner.style.width = "50px";
    spinner.style.height = "50px";
    spinner.style.border = "5px solid #f3f3f3";
    spinner.style.borderTop = "5px solid #3498db";
    spinner.style.borderRadius = "50%";
    spinner.style.animation = "spin 1s linear infinite";

    // Añadir estilo de animación
    const styleSheet = document.createElement("style");
    styleSheet.textContent = `
			@keyframes spin {
				0% { transform: rotate(0deg); }
				100% { transform: rotate(360deg); }
			}
		`;
    document.head.appendChild(styleSheet);

    overlay.appendChild(spinner);
    document.body.appendChild(overlay);
  }
}

function restablecerPantalla() {
  const overlay = document.getElementById("loading-overlay");
  if (overlay) {
    overlay.remove();
  }
}
/**
 * Garantiza la existencia de un contenedor de notificaciones en el documento
 *
 * @function ensureNotificationContainer
 * @description
 * - Verifica si existe un contenedor de notificaciones
 * - Crea el contenedor si no está presente
 * - Agrega el contenedor al cuerpo del documento
 *
 * @returns {HTMLDivElement} Contenedor de notificaciones
 *
 * @example
 * const container = ensureNotificationContainer();
 * // Contenedor de notificaciones garantizado
 */
function ensureNotificationContainer() {
  // Verificar si el contenedor ya existe
  let notificationContainer = document.getElementById("notificationContainer");

  // Si no existe, crearlo y agregarlo al body
  if (!notificationContainer) {
    notificationContainer = document.createElement("div");
    notificationContainer.id = "notificationContainer";
    notificationContainer.className = "notification-container";

    // Agregar el contenedor al body
    document.body.appendChild(notificationContainer);
  }

  return notificationContainer;
}

/**
 * Garantiza la existencia de un contenedor modal en el documento
 *
 * @function ensureModalContainer
 * @description
 * - Verifica si existe un modal
 * - Crea el modal con estructura predefinida si no está presente
 * - Incluye botón de cierre y contenedor de datos
 *
 * @returns {HTMLDivElement} Contenedor modal
 *
 * @example
 * const modal = ensureModalContainer();
 * // Modal garantizado con estructura completa
 */
function ensureModalContainer() {
  let modal = document.getElementById("miModal");

  if (!modal) {
    modal = document.createElement("div");
    modal.id = "miModal";
    modal.className = "modal";

    const modalContenido = document.createElement("div");
    modalContenido.className = "modal-contenido";

    const cerrarSpan = document.createElement("span");
    cerrarSpan.className = "cerrar-modal";
    cerrarSpan.innerHTML = "&times;";
    cerrarSpan.onclick = closeModal;

    const datosContainer = document.createElement("div");
    datosContainer.id = "datos";

    modalContenido.appendChild(cerrarSpan);
    modalContenido.appendChild(datosContainer);
    modal.appendChild(modalContenido);

    document.body.appendChild(modal);
  }

  return modal;
}

/**
 * Envía una solicitud AJAX con datos codificados en URL y maneja las respuestas y errores
 *
 * @param {string} method - Método HTTP de la solicitud (GET, POST, etc.)
 * @param {string} route - URL del endpoint del servidor
 * @param {Object|string} formUrlEncodedData - Datos a enviar en la solicitud
 * @param {string} successResponse - Valor esperado para considerar la solicitud exitosa
 * @param {string} successMessage - Mensaje de notificación para una respuesta exitosa
 * @param {string} unsuccessMessage - Mensaje de notificación para una respuesta no exitosa
 *
 * @description
 * - Oscurece la pantalla antes de enviar la solicitud
 * - Envía datos codificados en URL al servidor
 * - Muestra notificaciones según el resultado
 * - Maneja diferentes códigos de error HTTP con mensajes personalizados
 * - Registra errores detallados en la consola
 * - Restablece la pantalla después de completar la solicitud
 *
 * @example
 * sendUrlEncodedRequest(
 *   'POST',
 *   '/guardar-usuario',
 *   { nombre: 'Juan', email: 'juan@ejemplo.com' },
 *   'OK',
 *   'Usuario guardado exitosamente',
 *   'Error al guardar usuario'
 * )
 */
function sendUrlEncodedRequest(
  method,
  route,
  formUrlEncodedData,
  successResponse,
  successMessage,
  unsuccessMessage,
) {
  oscurecerPantalla();
  $.ajax({
    type: method,
    url: route,
    data: formUrlEncodedData,
    success: (data) => {
      if (data == successResponse) {
        showNotification(successMessage, "success");
      } else {
        showNotification(unsuccessMessage, "danger");
      }
      restablecerPantalla();
    },
    error: (xhr, status, error) => {
      // Manejo detallado de errores
      let errorMessage = "Error desconocido";

      // Diferentes tipos de manejo de error según el código de estado
      switch (xhr.status) {
        case 400:
          errorMessage = "Error de solicitud: Datos inválidos";
          break;
        case 401:
          errorMessage = "No autorizado: Inicie sesión nuevamente";
          break;
        case 403:
          errorMessage = "Acceso prohibido";
          break;
        case 404:
          errorMessage = "Recurso no encontrado";
          break;
        case 500:
          errorMessage = "Error interno del servidor";
          break;
        case 503:
          errorMessage = "Servicio no disponible";
          break;
      }

      // Mostrar mensaje de error
      $(target).html(`
						<div class="alert alert-danger">
							<strong>Error:</strong> ${errorMessage}
							<br>
							<small>Código de estado: ${xhr.status}</small>
						</div>
					`);

      // Log de error detallado para depuración
      console.error("Error en la solicitud:", {
        status: xhr.status,
        statusText: xhr.statusText,
        responseText: xhr.responseText,
        errorThrown: error,
      });
      restablecerPantalla();
    },
  });
}

/* LEGACY CODE
function sendMultiPartData(method, route, formId, successResponse, successMessage) {
  oscurecerPantalla();
  let formData = new FormData(document.getElementById(formId));
  $.ajax({
    type: method,
    url: route,
    data: formData,
    processData: false,
    contentType: false,
    success: (data) => {
      if (data == successResponse) {
        showNotification(successMessage, "success");
      } else {
        showNotification(data, "danger");
      }
      restablecerPantalla();
    },
    error: (xhr, status, error) => {
      // Manejo detallado de errores
      let errorMessage = "Error desconocido";

      // Diferentes tipos de manejo de error según el código de estado
      switch (xhr.status) {
        case 400:
          errorMessage = "Error de solicitud: Datos inválidos";
          break;
        case 401:
          errorMessage = "No autorizado: Inicie sesión nuevamente";
          break;
        case 403:
          errorMessage = "Acceso prohibido";
          break;
        case 404:
          errorMessage = "Recurso no encontrado";
          break;
        case 500:
          errorMessage = "Error interno del servidor";
          break;
        case 503:
          errorMessage = "Servicio no disponible";
          break;
      }

      // Mostrar mensaje de error
      $(target).html(`
						<div class="alert alert-danger">
							<strong>Error:</strong> ${errorMessage}
							<br>
							<small>Código de estado: ${xhr.status}</small>
						</div>
					`);

      // Log de error detallado para depuración
      console.error("Error en la solicitud:", {
        status: xhr.status,
        statusText: xhr.statusText,
        responseText: xhr.responseText,
        errorThrown: error,
      });
      restablecerPantalla();
    },
  });
}
  */

function sendUrlEncodedp(
  method,
  route,
  formUrlEncodedData,
  successResponse,
  successMessage,
) {
  oscurecerPantalla();
  $.ajax({
    type: method,
    url: route,
    data: formUrlEncodedData,
    success: (data) => {
      if (data == successResponse) {
        showNotification(successMessage, "success");
      } else {
        showNotification(data, "danger");
      }
      restablecerPantalla();
    },
    error: (xhr, status, error) => {
      // Manejo detallado de errores
      let errorMessage = "Error desconocido";

      // Diferentes tipos de manejo de error según el código de estado
      switch (xhr.status) {
        case 400:
          errorMessage = "Error de solicitud: Datos inválidos";
          break;
        case 401:
          errorMessage = "No autorizado: Inicie sesión nuevamente";
          break;
        case 403:
          errorMessage = "Acceso prohibido";
          break;
        case 404:
          errorMessage = "Recurso no encontrado";
          break;
        case 500:
          errorMessage = "Error interno del servidor";
          break;
        case 503:
          errorMessage = "Servicio no disponible";
          break;
      }

      // Mostrar mensaje de error
      $(target).html(`
						<div class="alert alert-danger">
							<strong>Error:</strong> ${errorMessage}
							<br>
							<small>Código de estado: ${xhr.status}</small>
						</div>
					`);

      // Log de error detallado para depuración
      console.error("Error en la solicitud:", {
        status: xhr.status,
        statusText: xhr.statusText,
        responseText: xhr.responseText,
        errorThrown: error,
      });
      restablecerPantalla();
    },
  });
}

/**
 * Carga contenido dinámicamente en un elemento objetivo mediante una solicitud AJAX
 *
 * @param {string} method - Método HTTP de la solicitud (GET, POST, etc.)
 * @param {string} route - URL del endpoint del servidor
 * @param {string} target - Selector CSS del elemento donde se cargará el contenido
 * @param {Object|string} formUrlEncodedData - Datos a enviar en la solicitud
 *
 * @description
 * - Oscurece la pantalla antes de enviar la solicitud
 * - Realiza una solicitud AJAX al servidor
 * - Inserta la respuesta directamente en el elemento objetivo
 * - Maneja diferentes códigos de error HTTP con mensajes personalizados
 * - Registra errores detallados en la consola
 * - Restablece la pantalla después de completar la solicitud
 *
 * @example
 * loadContentToTarget(
 *   'GET',
 *   '/obtener-usuarios',
 *   '#lista-usuarios',
 *   { pagina: 1 }
 * )
 */
function loadContentToTarget(method, route, target, formUrlEncodedData) {
  oscurecerPantalla();
  $.ajax({
    type: method,
    url: route,
    data: formUrlEncodedData,
    success: (data) => {
      $(target).html(data);
      restablecerPantalla();
    },
    error: (xhr, status, error) => {
      // Manejo detallado de errores
      let errorMessage = "Error desconocido";

      // Diferentes tipos de manejo de error según el código de estado
      switch (xhr.status) {
        case 400:
          errorMessage = "Error de solicitud: Datos inválidos";
          break;
        case 401:
          errorMessage = "No autorizado: Inicie sesión nuevamente";
          break;
        case 403:
          errorMessage = "Acceso prohibido";
          break;
        case 404:
          errorMessage = "Recurso no encontrado";
          break;
        case 500:
          errorMessage = "Error interno del servidor";
          break;
        case 503:
          errorMessage = "Servicio no disponible";
          break;
      }

      // Mostrar mensaje de error
      $(target).html(`
						<div class="alert alert-danger">
							<strong>Error:</strong> ${errorMessage}
							<br>
							<small>Código de estado: ${xhr.status}</small>
						</div>
					`);

      // Log de error detallado para depuración
      console.error("Error en la solicitud:", {
        status: xhr.status,
        statusText: xhr.statusText,
        responseText: xhr.responseText,
        errorThrown: error,
      });
      restablecerPantalla();
    },
  });
}

/**
 * Carga contenido dinámicamente con un spinner de carga y manejo de errores
 *
 * @param {string} method - Método HTTP de la solicitud (GET, POST, etc.)
 * @param {string} route - URL del endpoint del servidor
 * @param {string} target - Selector CSS del elemento donde se cargará el contenido
 * @param {Object|string} formUrlEncodedData - Datos a enviar en la solicitud
 *
 * @description
 * - Muestra un spinner de carga animado antes de la solicitud
 * - Realiza una solicitud AJAX al servidor
 * - Reemplaza el contenedor objetivo con la respuesta recibida
 * - Muestra un mensaje de error personalizado en caso de fallo
 * - Registra errores detallados en la consola
 * - No oscurece la pantalla completamente (comentario de oscurecerPantalla() desactivado)
 *
 * @example
 * LoadContents(
 *   'GET',
 *   '/obtener-contenido',
 *   '#contenedor-dinamico',
 *   { parametro: 'valor' }
 * )
 */
function LoadContents(method, route, target, formUrlEncodedData) {
  //oscurecerPantalla();
  $(target)
    .html(`<div class="centrado1"><div class="alert alert-warning" style="width: 50%"><img 
					src="/static/spint.gif" 
					alt="Cargando..." 
					style="
						width: 50px;
						height: 50px;
						margin-bottom: 20px;
						border-radius: 50%;
						box-shadow: 0 4px 6px rgba(0,0,0,0.1);
					"
				></div></div>`);
  $.ajax({
    type: method,
    url: route,
    data: formUrlEncodedData,
    success: (data) => {
      $(target).html(data);
    },
    error: (xhr, status, error) => {
      // Manejo detallado de errores
      let errorMessage = "Error desconocido";

      // Diferentes tipos de manejo de error según el código de estado
      switch (xhr.status) {
        case 400:
          errorMessage = "Error de solicitud: Datos inválidos";
          break;
        case 401:
          errorMessage = "No autorizado: Inicie sesión nuevamente";
          break;
        case 403:
          errorMessage = "Acceso prohibido";
          break;
        case 404:
          errorMessage = "Recurso no encontrado";
          break;
        case 500:
          errorMessage = "Error interno del servidor";
          break;
        case 503:
          errorMessage = "Servicio no disponible";
          break;
      }

      // Mostrar mensaje de error
      $(target).html(`
						<div class="alert alert-danger">
							<strong>Error:</strong> ${errorMessage}
							<br>
							<small>Código de estado: ${xhr.status}</small>
						</div>
					`);

      // Log de error detallado para depuración
      console.error("Error en la solicitud:", {
        status: xhr.status,
        statusText: xhr.statusText,
        responseText: xhr.responseText,
        errorThrown: error,
      });
    },
  });
}

var timerId;
/**
 * Envía una solicitud con un retraso para evitar múltiples solicitudes rápidas
 *
 * @param {string} method - Método HTTP de la solicitud (GET, POST, etc.)
 * @param {string} route - URL del endpoint del servidor
 * @param {string} target - Selector CSS del elemento donde se cargará el contenido
 * @param {string|*} e - Valor del parámetro a enviar
 * @param {string} name - Nombre del parámetro en la solicitud
 * @param {number} delay - Tiempo de espera en milisegundos antes de enviar la solicitud
 *
 * @description
 * - Cancela cualquier solicitud de temporizador pendiente
 * - Establece un nuevo temporizador con el retraso especificado
 * - Crea datos codificados en URL con el nombre y valor del parámetro
 * - Llama a LoadContents después del retraso para cargar contenido
 * - Útil para implementar búsqueda con retardo (debounce) o autocompletado
 *
 * @example
 * // Ejemplo de uso en un campo de búsqueda
 * sendRequestDelay(
 *   'GET',
 *   '/buscar',
 *   '#resultados-busqueda',
 *   inputValue,
 *   'termino',
 *   300
 * )
 */
function sendRequestDelay(method, route, target, e, name, delay) {
  clearTimeout(timerId);
  timerId = setTimeout(() => {
    let formUrlEncodedData = `${name}=${e}`;
    LoadContents(method, route, target, formUrlEncodedData);
  }, delay);
}

ensureNotificationContainer();
ensureModalContainer();
