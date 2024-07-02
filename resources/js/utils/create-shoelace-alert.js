/**
 * Creates a Shoelace alert with the given options.
 *
 * @param {Object} options - The options for the alert
 * @param {string} options.heading - The heading message for the alert
 * @param {"primary"|"success"|"neutral"|"danger"|"warning"} options.variant - The type of the alert
 * @param {number} options.duration - Duration for which the alert will be displayed
 * @param {string[]} options.messages - messages for the alert, every entry of message is in new line
 * @param {boolean} options.closable - is alert closable
 */
export function createShoelaceAlert(options) {
  const container = document.getElementById("alert-toast-wrapper");
  let icon;

  // Select icon based on variant
  switch (options.variant) {
    case "primary":
      icon = "info-circle";
      break;
    case "success":
      icon = "check2-circle";
      break;
    case "neutral":
      icon = "gear";
      break;
    case "warning":
      icon = "exclamation-triangle";
      break;
    case "danger":
      icon = "exclamation-octagon";
      break;
    default:
      icon = "info-circle";
  }

  let message = `<sl-icon name="${icon}" slot="icon"></sl-icon>`;

  if (options.heading) {
    message += `<p class="font-bold">${options.heading}</p>`;
  }

  if (options.messages) {
    options.messages.forEach((msg) => {
      message += `<p>${msg}</p>`;
    });
  }

  const alert = Object.assign(document.createElement("sl-alert"), {
    variant: options.variant,
    closable: options.closable,
    duration: options.duration,
    innerHTML: message,
  });

  // Append the alert to the alert-toast-wrapper container
  container.append(alert);
  alert.toast();
}
