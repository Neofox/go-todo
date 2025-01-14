import React, { useState, type ReactNode } from "react";

export interface ButtonProps {
    id: string;
    label?: string;
}

export function Button({ label, id }: ButtonProps): ReactNode {
    const [count, setCount] = useState(0);
    console.log("id coming from the server: ", id);

    return (
        <button className="bg-blue-500 text-white px-4 py-2 rounded-md" onClick={() => setCount(count + 1)}>
            {label && count === 0 ? label : `clicked: ${count}`}
        </button>
    );
}
