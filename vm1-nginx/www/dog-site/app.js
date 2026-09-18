const button = document.querySelector("#generate");
const dog = document.querySelector("#dog");
const name = document.querySelector("#name");
const description = document.querySelector("#description");
const error = document.querySelector("#error");

button.addEventListener("click", async () => {
  button.disabled = true;
  error.hidden = true;

  try {
    const response = await fetch("/api/dog");

    if (!response.ok) {
      throw new Error("Could not load a dog");
    }

    const data = await response.json();

    name.textContent = data.name;
    description.textContent = data.description;
    dog.hidden = false;
  } catch {
    dog.hidden = true;
    error.textContent = "Could not load a dog. Try again.";
    error.hidden = false;
  } finally {
    button.disabled = false;
  }
});
