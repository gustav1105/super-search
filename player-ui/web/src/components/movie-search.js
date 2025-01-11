import {createCard} from './movie-card.js';

export function createSearchForm() {
  const form = document.createElement("form");
  form.id = "searchForm";
  form.setAttribute('tabindex', '0');

  const queryLabel = document.createElement("label");
  queryLabel.setAttribute("for", "query");
  queryLabel.textContent = "Search Query";
  form.appendChild(queryLabel);

  const queryInput = document.createElement("input");
  queryInput.type = "text";
  queryInput.id = "query";
  queryInput.name = "query";
  queryInput.placeholder = "Enter search term";
  queryInput.required = true;
  form.appendChild(queryInput);
  form.appendChild(document.createElement("br"));

  const propertyLabel = document.createElement("label");
  propertyLabel.setAttribute("for", "property");
  propertyLabel.textContent = "Search Property";
  form.appendChild(propertyLabel);

  const propertySelect = document.createElement("select");
  propertySelect.id = "property";
  propertySelect.name = "property";
  propertySelect.required = true;

  const options = [
    { value: "title", text: "Title" },
    { value: "cast", text: "Cast" },
    { value: "genre", text: "Genre" },
    { value: "director", text: "Director" },
    { value: "release_date", text: "Release Date" },
    { value: "plot", text: "Plot" },
  ];

  options.forEach(optionData => {
    const option = document.createElement("option");
    option.value = optionData.value;
    option.textContent = optionData.text;
    propertySelect.appendChild(option);
  });

  form.appendChild(propertySelect);
  form.appendChild(document.createElement("br"));

  const submitButton = document.createElement("button");
  submitButton.type = "submit";
  submitButton.textContent = "Search";
  form.appendChild(submitButton);

  return form;
}

export function submitSearchForm (gridContainer) {
  const searchForm = document.getElementById("searchForm");
  searchForm.addEventListener("submit", async function (event) {
    event.preventDefault();

    const errorContainer = document.getElementById("errorContainer")
    const queryInput = document.getElementById("query");
    const propertySelect = document.getElementById("property");
    const query = queryInput.value;
    const property = propertySelect.value;

    // Clear gridContainer but keep searchForm intact

    gridContainer.replaceChildren(searchForm);

    // Clear the grid container before fetching new results
    try {
      // Send POST request to the server with user query and additional data
      const response = await fetch('http://192.168.10.94:8000/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          query: query,
          top_k: getValueByOption(property),
          property: property,
        }),
      });

      // Handle response errors
      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      // Parse the JSON response
      const data = await response.json();

      // Render results or display a no-results message
      if (data.results && data.results.length > 0) {
        data.results.forEach((item, i) => {
          const movie = item.metadata;
          const card = createCard(movie);
          card.setAttribute('data-index', i);
          card.setAttribute('tabindex', 0);
          gridContainer.appendChild(card);
        });
        const firstCard = gridContainer.querySelector('[data-index="0"]');
        if (firstCard) {
          firstCard.focus();
        }
      } else {
        resultContainer.innerHTML = '<p>No results found.</p>';
      }

    } catch (error) {
      // Display error messages
      errorContainer.textContent = `Error: ${error.message}`;
    }
  });
}

function getValueByOption(option) {
  switch (option) {
    case 'release_date':
      return 2000;
    case 'genre':
      return 5000;
    case 'title':
      return 200;
    case 'plot':
      return 500;
    case 'cast':
      return 300;
    default:
      return 500; // Default value if the option doesn't match
  }
}
