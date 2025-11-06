frontend/
├── index.html
├── package.json
├── postcss.config.js
├── tailwind.config.js
├── vite.config.js
│
└── src/
    ├── assets/                 # images, icons, logo, etc.
    │   └── logo.svg
    │
    ├── components/             # reusable UI building blocks
    │   ├── Sidebar.jsx
    │   ├── Navbar.jsx
    │   ├── CanvasArea.jsx
    │   ├── NodeCard.jsx
    │   └── Loader.jsx
    │
    ├── pages/                  # main app pages
    │   ├── Dashboard.jsx
    │   └── WorkflowBuilder.jsx
    │
    ├── context/                # global context or state providers
    │   └── WorkflowContext.jsx
    │
    ├── hooks/                  # custom React hooks
    │   └── useLocalStorage.js
    │
    ├── utils/                  # helper functions
    │   └── constants.js
    │
    ├── styles/                 # global/custom CSS (if needed)
    │   └── global.css
    │
    ├── App.jsx                 # root component
    ├── main.jsx                # app entry
    └── index.css               # Tailwind imports


Next step

A. Add node dragging inside the canvas (move nodes with mouse).
B. Add connectors (SVG lines) between nodes.
C. Build the properties panel to edit node text/props.
D. Replace the simple drag/drop with react-dnd for better UX.