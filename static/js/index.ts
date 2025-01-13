import React from "react";
import { createRoot } from "react-dom/client";
import { Button } from "./components/button";
import htmx from "htmx.org";

import "../css/input.css";

// Define window augmentation for htmx
declare global {
    interface Window {
        htmx: typeof htmx;
    }
}
window.htmx = htmx;

const componentRegistry = {
    Button: Button,
    // Add more components here as needed
    // 'Modal': Modal,
    // 'Card': Card,
};

type ComponentRegistry = typeof componentRegistry;

// Store references to root instances for cleanup
const reactRoots = new WeakMap<Element, ReturnType<typeof createRoot>>();

// Initialize React component for a single element
function initializeReactComponent(element: Element) {
    if (!reactRoots.has(element)) {
        const root = createRoot(element);
        reactRoots.set(element, root);
        const props = JSON.parse(element.getAttribute("data-react-props") || "{}");
        const componentName = element.getAttribute("data-react-component");

        if (!componentName || !(componentName in componentRegistry)) {
            console.warn(`Unknown or missing React component: ${componentName}`);
            return;
        }

        const Component = componentRegistry[componentName as keyof ComponentRegistry];
        root.render(React.createElement(Component, { ...props }));
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

// Initialize observer
const observer = new MutationObserver(mutations => {
    mutations.forEach(mutation => {
        // Handle removed nodes
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

// Initialize existing components
document.querySelectorAll("[data-react-component]").forEach(initializeReactComponent);
