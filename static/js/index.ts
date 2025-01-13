import React from "react";
import { createRoot } from "react-dom/client";
import htmx from "htmx.org";

import "../css/input.css";

// Define window augmentation for htmx
declare global {
    interface Window {
        htmx: typeof htmx;
    }
}
window.htmx = htmx;

// Store references to root instances for cleanup
const reactRoots = new WeakMap<Element, ReturnType<typeof createRoot>>();

// Initialize React component for a single element
async function initializeReactComponent(element: Element) {
    const componentName = element.getAttribute("data-react-component");

    // Skip if no component name or already initialized
    if (!componentName || reactRoots.has(element)) {
        return;
    }

    try {
        const root = createRoot(element);
        reactRoots.set(element, root);

        const props = JSON.parse(element.getAttribute("data-react-props") || "{}");

        // Dynamic import with explicit chunk name
        const Component = await import(
            /* webpackChunkName: "[request]" */
            `./components/${componentName.toLowerCase()}`
        ).then(m => m[componentName]);

        root.render(React.createElement(Component, { ...props }));
    } catch (error) {
        console.error(`Failed to load component ${componentName}:`, error);
    }
}

// Cleanup React component
function cleanupReactComponent(element: Element) {
    const root = reactRoots.get(element);
    if (root) {
        root.unmount();
        reactRoots.delete(element);
    }
}

// Initialize existing components at page load
document.querySelectorAll("[data-react-component]").forEach(initializeReactComponent);

// Initialize observer to look at new nodes being added to the DOM (with HTMX)
const observer = new MutationObserver(mutations => {
    mutations.forEach(mutation => {
        // Handle removed nodes (cleanup)
        // TODO: can probably be optimized, even removed
        mutation.removedNodes.forEach(node => {
            if (node instanceof Element) {
                const reactElements = [
                    ...(node.matches("[data-react-component]") ? [node] : []),
                    ...Array.from(node.querySelectorAll("[data-react-component]")),
                ];
                reactElements.forEach(cleanupReactComponent);
            }
        });

        // Handle added nodes
        mutation.addedNodes.forEach(node => {
            if (node instanceof Element) {
                const reactElements = [
                    ...(node.matches("[data-react-component]") ? [node] : []),
                    ...Array.from(node.querySelectorAll("[data-react-component]")),
                ];
                reactElements.forEach(initializeReactComponent);
            }
        });
    });
});

// Start observing
observer.observe(document.body, {
    childList: true,
    subtree: true,
});
