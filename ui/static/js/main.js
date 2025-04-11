/**
 * TechInsights Blog - Main JavaScript
 * Adds interactive features and animations to the blog
 */

document.addEventListener('DOMContentLoaded', function() {
  // Mobile menu toggle
  initMobileMenu();

  // Smooth scrolling for anchor links
  initSmoothScroll();

  // Initialize code syntax highlighting if prism.js is available
  if (typeof Prism !== 'undefined') {
    Prism.highlightAll();
  }

  // Add animation to cards
  initCardAnimations();

  // Initialize search functionality
  initSearch();
});

/**
 * Mobile menu functionality
 */
function initMobileMenu() {
  const menuToggle = document.getElementById('menu-toggle');
  const mobileMenu = document.getElementById('mobile-menu');

  if (menuToggle && mobileMenu) {
    menuToggle.addEventListener('click', function() {
      mobileMenu.classList.toggle('hidden');

      // Add slide animation
      if (!mobileMenu.classList.contains('hidden')) {
        mobileMenu.style.maxHeight = mobileMenu.scrollHeight + 'px';
      } else {
        mobileMenu.style.maxHeight = '0';
      }
    });
  }
}

/**
 * Smooth scrolling for anchor links
 */
function initSmoothScroll() {
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
      const targetId = this.getAttribute('href');

      if (targetId === '#') return;

      e.preventDefault();

      const targetElement = document.querySelector(targetId);

      if (targetElement) {
        window.scrollTo({
          top: targetElement.offsetTop - 80, // Offset for fixed header
          behavior: 'smooth'
        });
      }
    });
  });
}

/**
 * Card animations and hover effects
 */
function initCardAnimations() {
  const cards = document.querySelectorAll('.card-hover');

  cards.forEach(card => {
    // Add subtle entrance animation
    card.style.opacity = '0';
    card.style.transform = 'translateY(20px)';

    setTimeout(() => {
      card.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
      card.style.opacity = '1';
      card.style.transform = 'translateY(0)';
    }, 100);
  });
}

/**
 * Search functionality
 */
function initSearch() {
  const searchForms = document.querySelectorAll('.search-form');

  searchForms.forEach(form => {
    form.addEventListener('submit', function(e) {
      e.preventDefault();

      const searchInput = this.querySelector('input[type="search"]');
      const searchTerm = searchInput.value.trim();

      if (searchTerm.length < 2) {
        // Show error message for short search terms
        const errorMsg = document.createElement('div');
        errorMsg.className = 'text-red-500 text-sm mt-1';
        errorMsg.textContent = 'Please enter at least 2 characters';

        // Remove any existing error message
        const existingError = form.querySelector('.text-red-500');
        if (existingError) {
          existingError.remove();
        }

        form.appendChild(errorMsg);

        // Auto-remove error after 3 seconds
        setTimeout(() => {
          errorMsg.remove();
        }, 3000);

        return;
      }

      // Redirect to search results page (this would be implemented server-side)
      window.location.href = `/search?q=${encodeURIComponent(searchTerm)}`;
    });
  });
}

/**
 * Theme toggle functionality (if implemented)
 */
function initThemeToggle() {
  const themeToggle = document.getElementById('theme-toggle');

  if (themeToggle) {
    themeToggle.addEventListener('click', function() {
      document.documentElement.classList.toggle('dark-mode');

      // Save preference to localStorage
      const isDarkMode = document.documentElement.classList.contains('dark-mode');
      localStorage.setItem('darkMode', isDarkMode ? 'true' : 'false');
    });

    // Check for saved preference
    const savedDarkMode = localStorage.getItem('darkMode') === 'true';
    if (savedDarkMode) {
      document.documentElement.classList.add('dark-mode');
    }
  }
}

/**
 * Newsletter form validation
 */
document.addEventListener('DOMContentLoaded', function() {
  const newsletterForms = document.querySelectorAll('.newsletter-form');

  newsletterForms.forEach(form => {
    form.addEventListener('submit', function(e) {
      e.preventDefault();

      const emailInput = this.querySelector('input[type="email"]');
      const email = emailInput.value.trim();

      // Simple email validation
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

      if (!emailRegex.test(email)) {
        // Show error for invalid email
        emailInput.classList.add('border-red-500');

        const errorMsg = document.createElement('div');
        errorMsg.className = 'text-red-100 text-sm mt-1';
        errorMsg.textContent = 'Please enter a valid email address';

        // Remove any existing error message
        const existingError = form.querySelector('.text-red-100');
        if (existingError) {
          existingError.remove();
        }

        form.appendChild(errorMsg);
        return;
      }

      // Simulate form submission
      const submitButton = form.querySelector('button[type="submit"]');
      const originalText = submitButton.textContent;

      submitButton.disabled = true;
      submitButton.textContent = 'Subscribing...';

      // Simulate API call
      setTimeout(() => {
        // Success state
        form.innerHTML = `
          <div class="text-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-white mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <p class="text-white font-medium">Thanks for subscribing!</p>
            <p class="text-white opacity-80 text-sm mt-1">We'll keep you updated with the latest tech news.</p>
          </div>
        `;
      }, 1500);
    });
  });
});