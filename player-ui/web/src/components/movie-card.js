export function createStarRating(rating) {
    const roundedRating = Math.round(rating); // Round to nearest integer
    const starsContainer = document.createElement('div');
    starsContainer.className = "stars-container";

    for (let i = 0; i < 10; i++) {
        const star = document.createElement('span');
        star.textContent = i < roundedRating ? '★' : '☆'; // Filled or empty star
        starsContainer.appendChild(star);
    }

    return starsContainer;
}

function createCardContent(movie) {
  const content = document.createElement('div');
  content.className = 'card-content';

  const title = createInfoElement('h3', '', movie.title || 'No Title');
  const rating = createStarRating(movie.rating || 0);
  const cast = createInfoElement('p', 'Cast:', movie.cast);
  const director = createInfoElement('p', 'Director:', movie.director);
  const genre = createInfoElement('p', 'Genre:', movie.genre);
  const releaseDate = createInfoElement('p', 'Release Date:', movie.release_date);
  const plot = createInfoElement('p', 'Plot:', movie.plot);
  content.appendChild(title);
  content.appendChild(rating);
  cast.style.display = 'none';

  cast.classList.add('cast');
  content.appendChild(cast);
  //director.style.display = 'none';
  content.appendChild(director);
  genre.style.display = 'none';

  genre.classList.add('genre');
  content.appendChild(genre);
  //releaseDate.style.display = 'none';
  content.appendChild(releaseDate);
  plot.style.display = 'none';
  plot.classList.add('plot')
  content.appendChild(plot);

  const menu = createCardMenuLayout(movie.tmdb_id);
  content.appendChild(menu);
  return content;
}

function createYoutubeFrame(youtubeId) {
  const iframeContainer = document.createElement('div');
  iframeContainer.className = 'iframe-container';

  const iframe = document.createElement('iframe');
  iframe.className = 'iframe';


  iframe.id = youtubeId ? `https://www.youtube-nocookie.com/embed/${youtubeId}?enablejsapi=1` : '';
  iframe.frameBorder = '0';
  iframe.allowFullscreen = true;
  iframeContainer.appendChild(iframe);
  iframeContainer.style.display = 'none';
  iframeContainer.classList.add('youtube');
  iframeContainer.setAttribute("tabindex", '0');
  return iframeContainer;
}

export function createCard(movie) {
  // Create the main card container
  const card = document.createElement('div');
  card.className = 'movie-card';

  const proxyImageUrl = `http://192.168.10.94:3000/proxy-image?url=${encodeURIComponent(movie.movie_image)}`;
  card.style.backgroundImage = `url('${proxyImageUrl}')`;
  card.dataset.url = `url('${proxyImageUrl}')`;
  const content = createCardContent(movie);
  card.appendChild(content);

  // Create the YouTube iframe section if YouTube ID exists
  const youtubeFrame = createYoutubeFrame(movie.youtube_trailer);
  card.appendChild(youtubeFrame);


  card.setAttribute('tabindex', '0');
  return card;
}

function createInfoElement(tag, label, value, className = '') {
  const element = document.createElement(tag);
  if (className) {
    element.className = className;
  }
  element.textContent = label ? `${label} ${value || 'N/A'}` : value || 'N/A';
  return element;
}

function createCardMenuLayout(tmdbId) {
  const menu = document.createElement('div');
  menu.classList.add("menu");

  const backToSearchButton = createBackToSearchFormButton();
  const playButton = createPlayButton();
  const infoButton = createInfoButton(tmdbId);
  menu.appendChild(backToSearchButton);
  menu.appendChild(infoButton);
  menu.appendChild(playButton);

  return menu;
}

function createPlayButton() {
  const button = document.createElement("button");
  button.id = "playButton";
  button.style.borderRadius = "50%";
  button.style.border = "none";
  button.style.cursor = "pointer";
  button.style.backgroundColor = "transparent";
  button.style.display ="flex";
  button.classList.add("play");
  button.setAttribute('tabindex', 0);

  const icon = document.createElement("i");
  icon.className = "material-icons";
  icon.textContent = "play_arrow";
  icon.style.fontSize ="42px";
  icon.style.color = "white";

  button.appendChild(icon); 
  button.addEventListener("click", () => {
  });

  return button; 
}

function createBackToSearchFormButton() {
  const button = document.createElement("button");
  button.id = "backToSearchFormButton";
  button.style.borderRadius = "50%";
  button.style.border = "none";
  button.style.cursor = "pointer";
  button.style.backgroundColor = "transparent";
  button.style.display ="flex";
  button.classList.add("back");
  button.setAttribute('tabindex', 0);
 
  const icon = document.createElement("i");
  icon.className = "material-icons";
  icon.textContent = "search";
  icon.style.fontSize ="42px";
  icon.style.color = "white";
  
  button.appendChild(icon); 
  button.addEventListener("click", () => {
    const searchInput = document.getElementById("searchForm");
    if (searchInput) {
      searchInput.focus(); // Focus the search input
    }
  });

  return button; // Return the button element
}

function createInfoButton(tmdbId) {
  const button = document.createElement("button");
  button.id = "infoButton";
  button.style.borderRadius = "50%";
  button.style.border = "none";
  button.style.cursor = "pointer";
  button.style.backgroundColor = "transparent";
  button.style.display ="flex";
  button.classList.add("back");
  button.setAttribute('tabindex', 0);
  button.dataset.tmdbId = tmdbId;
 
  const icon = document.createElement("i");
  icon.className = "material-icons";
  icon.textContent = "info";
  icon.style.fontSize ="42px";
  icon.style.color = "white";
  
  button.appendChild(icon); 
  button.addEventListener("click", () => {

  });

  return button; // Return the button element
}
