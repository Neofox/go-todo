import React, { useState, type ReactNode } from "react";

export function Button(props: ButtonProps): ReactNode {
    const [count, setCount] = useState(0);

    console.log(props);
    return (
        <button className="bg-blue-300 text-white px-4 py-2 rounded-md" onClick={() => setCount(count + 1)}>
            clicked: {count}
        </button>
    );
}

export interface ButtonProps {
    id: string;
}
