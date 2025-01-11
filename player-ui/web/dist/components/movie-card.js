export function createStarRating(rating) {
  var roundedRating = Math.round(rating); // Round to nearest integer
  var starsContainer = document.createElement('div');
  starsContainer.className = "stars-container";
  for (var i = 0; i < 10; i++) {
    var star = document.createElement('span');
    star.textContent = i < roundedRating ? '★' : '☆'; // Filled or empty star
    starsContainer.appendChild(star);
  }
  return starsContainer;
}
function createCardContent(movie) {
  var content = document.createElement('div');
  content.className = 'card-content';
  var title = createInfoElement('h3', '', movie.title || 'No Title');
  var rating = createStarRating(movie.rating || 0);
  var cast = createInfoElement('p', 'Cast:', movie.cast);
  var director = createInfoElement('p', 'Director:', movie.director);
  var genre = createInfoElement('p', 'Genre:', movie.genre);
  var releaseDate = createInfoElement('p', 'Release Date:', movie.release_date);
  var plot = createInfoElement('p', 'Plot:', movie.plot);
  var tmbd = createInfoElement('p', 'TMDB:', movie.tmbd_id);
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
  plot.classList.add('plot');
  content.appendChild(plot);
  tmbd.style.display = 'none';
  tmbd.classList.add('tmdb');
  content.appendChild(tmbd);
  return content;
}
function createYoutubeFrame(youtubeId) {
  var iframeContainer = document.createElement('div');
  iframeContainer.className = 'iframe-container';
  var iframe = document.createElement('iframe');
  iframe.className = 'iframe';
  iframe.id = youtubeId ? "https://www.youtube-nocookie.com/embed/".concat(youtubeId, "?enablejsapi=1") : '';
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
  var card = document.createElement('div');
  card.className = 'movie-card';
  var proxyImageUrl = "http://192.168.10.94:3000/proxy-image?url=".concat(encodeURIComponent(movie.movie_image));
  card.style.backgroundImage = "url('".concat(proxyImageUrl, "')");
  card.dataset.url = "url('".concat(proxyImageUrl, "')");
  var content = createCardContent(movie);
  card.appendChild(content);

  // Create the YouTube iframe section if YouTube ID exists
  var youtubeFrame = createYoutubeFrame(movie.youtube_trailer);
  card.appendChild(youtubeFrame);
  var backToSearchButton = createBackToSearchFormButton();
  card.appendChild(backToSearchButton);
  card.setAttribute('tabindex', '0');
  return card;
}
function createInfoElement(tag, label, value) {
  var className = arguments.length > 3 && arguments[3] !== undefined ? arguments[3] : '';
  var element = document.createElement(tag);
  if (className) {
    element.className = className;
  }
  element.textContent = label ? "".concat(label, " ").concat(value || 'N/A') : value || 'N/A';
  return element;
}
function createBackToSearchFormButton() {
  var button = document.createElement("button"); // Properly initialize the button
  button.textContent = "🔍"; // Set button text
  button.id = "backToSearchFormButton"; // Optional: Assign an ID
  button.style.width = "56px";
  button.style.height = "56px";
  button.style.borderRadius = "50%";
  button.style.border = "none";
  button.style.fontSize = "24px";
  button.style.cursor = "pointer";
  button.style.justifyContent = "center";
  button.style.alignItems = "center";
  button.style.position = "absolute";
  button.style.left = "10px";
  button.style.bottom = "20px";
  button.style.backgroundColor = "transparent";
  button.style.display = 'none';
  button.classList.add("back");

  // Add a click event listener
  button.addEventListener("click", function () {
    var searchInput = document.getElementById("searchForm");
    if (searchInput) {
      searchInput.focus(); // Focus the search input
    }
  });
  return button; // Return the button element
}