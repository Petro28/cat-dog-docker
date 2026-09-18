const button = document.querySelector("#generate");
const result = document.querySelector("#result");
const name = document.querySelector("#cat-name");
const description = document.querySelector("#cat-description");
const error = document.querySelector("#error");

button.addEventListener("click", async () => {
  error.textContent = "";
  button.disabled = true;

  try {
    const response = await fetch("/api/cat");

    if (!response.ok) {
      throw new Error("API request failed");
    }

    const cat = await response.json();

    name.textContent = cat.name;
    description.textContent = cat.description;
    result.hidden = false;
  } catch {
    error.textContent = "Could not load a cat. Try again.";
  } finally {
    button.disabled = false;
  }
});
