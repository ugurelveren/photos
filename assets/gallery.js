// Get elements
const gallery = document.getElementById('gallery');
const modal = document.getElementById('imageModal');
const modalImage = document.getElementById('modalImage');
const caption = document.getElementById('caption');
const prevButton = document.getElementById('prev');
const nextButton = document.getElementById('next');
const closeButton = document.querySelector('.close');

let images = [];
let currentIndex = 0;

// Fetch images from images.json
fetch('data/images.json')
  .then(response => response.json())
  .then(data => {
    images = data.images;
    loadGallery();
  })
  .catch(error => console.error('Error loading images:', error));

// Load images into the gallery
function loadGallery() {
  images.forEach((image, index) => {
    const imgElement = document.createElement('img');
    imgElement.src = image.url;
    imgElement.alt = "Creative Commons Attribution-NonCommercial 4.0 International Public License";
    imgElement.addEventListener('click', () => openModal(index));
    gallery.appendChild(imgElement);
  });
}

// Open the modal with the selected image
function openModal(index) {
  currentIndex = index;
  modal.style.display = 'block';
  modalImage.src = images[currentIndex].url;
  caption.textContent = "Creative Commons Attribution-NonCommercial 4.0 International Public License";
}

// Close the modal
closeButton.addEventListener('click', () => {
  modal.style.display = 'none';
});

// Next/Previous controls
nextButton.addEventListener('click', () => {
  currentIndex = (currentIndex + 1) % images.length;
  openModal(currentIndex);
});

prevButton.addEventListener('click', () => {
  currentIndex = (currentIndex - 1 + images.length) % images.length;
  openModal(currentIndex);
});

// Close the modal when clicking outside the image
window.addEventListener('click', (event) => {
  if (event.target === modal) {
    modal.style.display = 'none';
  }
});