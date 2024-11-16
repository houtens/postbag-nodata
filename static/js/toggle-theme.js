window.addEventListener('DOMContentLoaded', () => {
  const themeToggleButton = document.getElementById('theme-toggle');
  themeToggleButton.addEventListener('click', () => {
    const currentTheme = document.documentElement.getAttribute("data-theme");
    const newTheme = currentTheme === "light" ? "dark" : "light";
    document.documentElement.setAttribute('data-theme', newTheme)
  })
})
// const themeToggleButton = document.getElementById("theme-toggle");

// themeToggleButton.addEventListener("click", () => {
//   const currentTheme = document.documentElement.getAttribute("data-theme");
//   const newTheme = currentTheme === "light" ? "dark" : "light";
//   document.documentElement.setAttribute("data-theme", newTheme);
// });
