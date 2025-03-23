// DOM Elements
const markdownEditor = document.getElementById('markdown-editor');
const markdownPreview = document.getElementById('markdown-preview');
const decreaseFontBtn = document.getElementById('decrease-font');
const increaseFontBtn = document.getElementById('increase-font');
const toggleModeBtn = document.getElementById('toggle-mode');
const editorContainer = document.getElementById('editor-container');

// State Variables
let baseFontSize = 16;
let undoStack = [];
let redoStack = [];
let lastChange = '';
let isReaderMode = false;

// Initialize Application
document.addEventListener('DOMContentLoaded', function() {
  // Set initial content and update preview
  markdownEditor.value = '# Welcome to Markdown Scratchpad\n\nStart typing your markdown here...\n\n## Features\n\n- **Bold** and *italic* text\n- [Links](https://example.com)\n- Lists (ordered and unordered)\n- Code blocks\n- And more!\n\n```\nfunction example() {\n  console.log("Hello, Markdown!");\n}\n```';
  updatePreview();
  
  // Setup event listeners
  setupEventListeners();
  
  // Save initial state for undo
  saveState();
});

// Setup event listeners for the application
function setupEventListeners() {
  // Update preview when markdown changes
  markdownEditor.addEventListener('input', function() {
    updatePreview();
    saveState();
  });
  
  // Font size control
  decreaseFontBtn.addEventListener('click', decreaseFontSize);
  increaseFontBtn.addEventListener('click', increaseFontSize);
  
  // Handle tab key for indentation
  markdownEditor.addEventListener('keydown', function(e) {
    if (e.key === 'Tab') {
      e.preventDefault();
      const start = this.selectionStart;
      const end = this.selectionEnd;
      
      // Insert tab at position
      this.value = this.value.substring(0, start) + '  ' + this.value.substring(end);
      
      // Move cursor after the inserted tab
      this.selectionStart = this.selectionEnd = start + 2;
      updatePreview();
      saveState();
    }
  });
}

// Update the markdown preview with rendered HTML
function updatePreview() {
  // Configure marked options
  marked.setOptions({
    breaks: true,
    gfm: true,
    headerIds: true,
    highlight: function(code) {
      return code;
    }
  });
  
  // Render markdown to HTML
  markdownPreview.innerHTML = marked.parse(markdownEditor.value);
  
  // Make links open in new tab
  const links = markdownPreview.querySelectorAll('a');
  links.forEach(link => {
    link.setAttribute('target', '_blank');
    link.setAttribute('rel', 'noopener noreferrer');
  });
}

// Font size controls
function decreaseFontSize() {
  if (baseFontSize > 10) {
    baseFontSize -= 1;
    updateFontSize();
  }
}

function increaseFontSize() {
  if (baseFontSize < 30) {
    baseFontSize += 1;
    updateFontSize();
  }
}

function updateFontSize() {
  markdownEditor.style.fontSize = `${baseFontSize}px`;
  markdownPreview.style.fontSize = `${baseFontSize}px`;
  
  const headings = markdownPreview.querySelectorAll('h1, h2, h3, h4, h5, h6');
  headings.forEach(heading => {
    const level = parseInt(heading.tagName.substring(1));
    const scaleFactor = Math.max(1 + (1.5 - level * 0.1), 1);
    heading.style.fontSize = `${baseFontSize * scaleFactor}px`;
  });
  
  const codeElements = markdownPreview.querySelectorAll('code');
  codeElements.forEach(code => {
    code.style.fontSize = `${baseFontSize * 0.9}px`;
  });
}

// Select all text in editor
function selectAllText() {
  markdownEditor.select();
  markdownEditor.focus();
}

// Toggle between reader and writer mode
function toggleMode() {
  isReaderMode = !isReaderMode;
  editorContainer.classList.toggle('reader-mode', isReaderMode);
  
  // Update button icon and title
  if (isReaderMode) {
    toggleModeBtn.innerHTML = '<i class="fas fa-pencil-alt"></i>';
    toggleModeBtn.title = "Write Mode";
  } else {
    toggleModeBtn.innerHTML = '<i class="fas fa-book-reader"></i>';
    toggleModeBtn.title = "Reader Mode";
    markdownEditor.focus();
  }
  
  // Make sure preview is updated
  updatePreview();
}

// Undo/Redo functionality
function saveState() {
  const currentValue = markdownEditor.value;
  
  // Only save state if it's different from last saved
  if (currentValue !== lastChange) {
    undoStack.push(lastChange);
    lastChange = currentValue;
    
    // Clear redo stack when new changes are made
    redoStack = [];
    
    // Limit undo stack size
    if (undoStack.length > 100) {
      undoStack.shift();
    }
  }
}

function undoChange() {
  if (undoStack.length > 0) {
    const currentValue = markdownEditor.value;
    redoStack.push(currentValue);
    
    const previousValue = undoStack.pop();
    markdownEditor.value = previousValue;
    lastChange = previousValue;
    
    updatePreview();
  }
}

function redoChange() {
  if (redoStack.length > 0) {
    const currentValue = markdownEditor.value;
    undoStack.push(currentValue);
    
    const nextValue = redoStack.pop();
    markdownEditor.value = nextValue;
    lastChange = nextValue;
    
    updatePreview();
  }
}
