import { createSearchForm, submitSearchForm } from './components/movie-search.js';
import { createCard } from './components/movie-card.js';

document.addEventListener('DOMContentLoaded', () => {
  const app = document.getElementById('app');
  if (app) {
    // Error Container
    const errorContainer = document.createElement("div");
    errorContainer.id = "errorContainer";
    errorContainer.style.color = "red";
    errorContainer.style.marginTop = "8px";
    app.appendChild(errorContainer);

    // Create Grid Container
    const gridContainer = document.createElement("div");
    gridContainer.id = "gridContainer";
    gridContainer.style.maxHeight = `${window.innerHeight}px`;

    const searchForm = createSearchForm();
    app.appendChild(searchForm);
    app.appendChild(gridContainer);

    const inputField = searchForm.querySelector("input");
    if (inputField) {
      inputField.focus();
    }
    submitSearchForm(gridContainer);
  }

  document.addEventListener("keydown", (event) => {
    const focusedCard = document.activeElement;
    if (focusedCard.classList.contains("movie-card")) {
      const gridContainer = document.getElementById("gridContainer");
      const gridItems = Array.from(gridContainer.children);
      const focusedIndex = gridItems.indexOf(focusedCard);

      if (event.key === "ArrowRight") {
        moveFocus(focusedIndex + 1, gridItems);
      } else if (event.key === "ArrowLeft") {
        moveFocus(focusedIndex - 1, gridItems);
      } else if (event.key === "ArrowDown") {
        moveFocus(focusedIndex + 3, gridItems); // Move to the next row
      } else if (event.key === "ArrowUp") {
        moveFocus(focusedIndex - 3, gridItems); // Move to the previous row
      } else if (event.key === "Enter") {
        setCardFocus(focusedCard, gridContainer);
      }
    } else if (focusedCard.classList.contains("youtube")) {
      const iframe = focusedCard.querySelector("iframe");
      if(event.key === "Enter") {
        if (iframe) {
          iframe.contentWindow.postMessage(
            JSON.stringify({
              event: "command",
              func: "playVideo",
              args: []
            }),
            "*"
          );
        }  
      } else if (event.key === "ArrowDown" || event.key === "ArrowUp" || event.key === "ArrowLeft" || event.key === "ArrowRight") {
        if (iframe) {
          iframe.contentWindow.postMessage(
            JSON.stringify({
              event: "command",
              func: "stopVideo",
              args: []
            }),
            "*"
          );
        }  
      }
    }
  });
});

function moveFocus(targetIndex, gridItems) {
  if (targetIndex >= 0 && targetIndex < gridItems.length) {
    gridItems[targetIndex].focus();
  }
}

function setCardFocus(focusedCard, gridContainer) {
  const gridItems = Array.from(gridContainer.children);
  const focusedIndex = gridItems.indexOf(focusedCard);
  const rowStartIndex = Math.floor(focusedIndex / 3) * 3; // Start of the row

  // Expand the card to span all three columns
  focusedCard.style.gridColumn = "1 / span 3";
  focusedCard.style.transition = "all 0.3s ease";

  const plotElement = focusedCard.querySelector('.plot');
  if (plotElement) {
    plotElement.style.display = "block";
  }

  const genreElement = focusedCard.querySelector('.genre');
  if (genreElement) {
    genreElement.style.display = "block";
  }

  const castElement = focusedCard.querySelector('.cast');
  if (castElement) {
    castElement.style.display = "block";
  }


  gridContainer.insertBefore(focusedCard, gridItems[rowStartIndex]);

  const content = focusedCard.querySelector('.card-content');
  const youtubeTrailerElement = focusedCard.querySelector('.youtube');

  if(youtubeTrailerElement) {
    youtubeTrailerElement.style.display = 'block'
    const cardBackground = window.getComputedStyle(focusedCard).backgroundImage;
    youtubeTrailerElement.style.backgroundImage = cardBackground;
    youtubeTrailerElement.focus();

    const iframe = youtubeTrailerElement.querySelector('iframe');
    iframe.src = iframe.id;
    iframe.setAttribute('tabindex', 0);
    focusedCard.style.boxShadow  = "none";
    focusedCard.style.border = "none";
    focusedCard.style.borderRadius = 0;
  }

  const cardContent = focusedCard.querySelector('.card-content'); // Adjust the selector if needed
  if (cardContent) {
    cardContent.style.paddingTop = "10px";
    cardContent.style.marginTop = "0px";
    cardContent.style.backgroundColor = "transparent";
  }

  const cardMenu = focusedCard.querySelector('.menu');
  if(cardMenu) {
    cardMenu.style.display = "flex";
  }

  focusedCard.style.backgroundImage = "";
  // Reset other cards
  gridItems.forEach((item) => {
    if (item !== focusedCard) {
      const itemMenu = item.querySelector('.menu');
      if(itemMenu) {
        itemMenu.style.display = "none";
      }

      const itemContent = item.querySelector('.card-content');
      if(itemContent) {
        itemContent.style.marginTop = "auto";
        itemContent.style.backgroundColor = "rgba(0, 0, 0, 0.7)";
      }

      item.style.backgroundImage = item.dataset.url;
      item.style.gridColumn = "auto";
      item.style.boxShadow  = "0 4px 8px rgba(0, 0, 0, 0.1)";
      item.style.border = "1px solid #ccc";
      item.style.borderRadius = "8px";

      const itemCast = item.querySelector('.cast');
      if (itemCast) {
        itemCast.style.display = "none";
      }
      const itemGenre = item.querySelector('.genre');
      if (itemGenre) {
        itemGenre.style.display = "none";
      }

      const itemPlot = item.querySelector('.plot');
      if (itemPlot) {
        itemPlot.style.display = "none";
      }

      const itemYoutube = item.querySelector('.youtube');
      if (itemYoutube) {
        itemYoutube.style.display = "none";
      }
    }
  });
}

