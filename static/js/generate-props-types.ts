import * as path from "path";
import * as ts from "typescript";
import { readdir } from "node:fs/promises";

function getTypeFromTsType(type: ts.Type, typeChecker: ts.TypeChecker): string {
    // Handle array types
    if (type.symbol?.name === "Array") {
        const typeArgs = (type as ts.TypeReference).typeArguments;
        if (typeArgs && typeArgs.length > 0) {
            const elementType = getTypeFromTsType(typeArgs[0], typeChecker);
            return `[]${elementType}`;
        }
        return "[]any";
    }

    // Handle function types
    if (type.getCallSignatures().length > 0) {
        return handleFunctionType(type, typeChecker);
    }

    if (type.isStringLiteral() || typeChecker.typeToString(type) === "string") return "string";

    if (type.symbol?.name === "Array") {
        return handleArrayType(type, typeChecker);
    }

    // Handle object types
    if (type.symbol?.members && type.symbol.members.size > 0) {
        return handleObjectType(type, typeChecker);
    }

    if (type.getFlags() & ts.TypeFlags.Boolean) {
        return "bool";
    }

    // Handle number types
    if (type.isNumberLiteral() || typeChecker.typeToString(type) === "number")
        return handleNumberType(type, typeChecker);

    // Handle union types
    if (type.isUnion()) {
        // For now, we'll treat unions as any
        // Could be enhanced to handle specific cases
        return "any";
    }

    return "any";
}

function handleArrayType(type: ts.Type, typeChecker: ts.TypeChecker): string {
    const elementType = typeChecker.getTypeArguments(type as ts.TypeReference)[0];

    // Check if it's an object type with properties using proper type flags
    if (elementType.getProperties && elementType.getProperties().length > 0) {
        const properties = elementType.getProperties();
        const structFields = properties
            .map(prop => {
                const propType = typeChecker.getTypeOfSymbol(prop);
                const goType = getTypeFromTsType(propType, typeChecker);
                return `${capitalize(prop.name)} ${goType} \`json:"${prop.name}"\``;
            })
            .join("\n        ");

        return `[]struct {\n        ${structFields}\n    }`;
    }

    // For primitive array types
    const elementGoType = getTypeFromTsType(elementType, typeChecker);
    return `[]${elementGoType}`;
}

function handleObjectType(type: ts.Type, typeChecker: ts.TypeChecker): string {
    // Check if it's an object type with properties
    if (type.getProperties && type.getProperties().length > 0) {
        const properties = type.getProperties();
        const structFields = properties
            .map(prop => {
                const propType = typeChecker.getTypeOfSymbol(prop);
                const goType = getTypeFromTsType(propType, typeChecker);
                return `${capitalize(prop.name)} ${goType} \`json:"${prop.name}"\``;
            })
            .join("\n        ");

        return `struct {\n        ${structFields}\n    }`;
    }

    // Fallback to any if we can't determine the structure
    return "any";
}

function handleNumberType(type: ts.Type, typeChecker: ts.TypeChecker): string {
    // Handle regular number type
    if (type.flags & ts.TypeFlags.Number) return "float64";

    if (type.flags & ts.TypeFlags.NumberLiteral) {
        const typeText = typeChecker.typeToString(type);
        const value = parseFloat(typeText);
        return Number.isInteger(value) ? "int" : "float64";
    }

    const value = (type as ts.NumberLiteralType).value;
    return Number.isInteger(value) ? "int" : "float64";
}

function handleFunctionType(type: ts.Type, typeChecker: ts.TypeChecker): string {
    const signature = type.getCallSignatures()[0];
    const parameters = signature.getParameters();
    const returnType = signature.getReturnType();

    // Handle parameters
    const paramTypes = parameters.map(param => {
        const paramType = typeChecker.getTypeOfSymbol(param);
        return getTypeFromTsType(paramType, typeChecker);
    });

    // Handle return type
    const returnTypeStr = getTypeFromTsType(returnType, typeChecker);

    if (paramTypes.length === 0) {
        return returnTypeStr === "any" ? "func()" : `func() ${returnTypeStr}`;
    }

    return `func(${paramTypes.join(", ")}) ${returnTypeStr}`;
}

function generateGoType(interfaceName: string, properties: ts.Symbol[], typeChecker: ts.TypeChecker): string {
    const goProperties = properties
        .map(prop => {
            const name = prop.getName();
            const type = typeChecker.getTypeOfSymbol(prop);
            const declarations = prop.declarations;

            let goType = getTypeFromTsType(type, typeChecker);
            const isOptional = declarations?.some(d => ts.isPropertySignature(d) && d.questionToken !== undefined);
            const isFunction = type.getCallSignatures().length > 0;

            let jsonTag = "";
            if (isFunction) {
                jsonTag = 'json:"-"';
            } else {
                jsonTag = isOptional ? `json:"${name.toLowerCase()},omitempty"` : `json:"${name.toLowerCase()}"`;
            }

            return `    ${capitalize(name)} ${goType} \`${jsonTag}\``;
        })
        .join("\n");

    return `
// ${interfaceName} represents the props for the ${interfaceName.replace("Props", "")} component
type ${interfaceName} struct {
${goProperties}
}

func (p ${interfaceName}) String() (string, error) {
    b, err := json.Marshal(p)
    if err != nil {
        return "", err
    }
    return string(b), nil
}

`;
}

async function generateGoTypes() {
    const componentsDir = "./static/js/components";

    const files = (await readdir(componentsDir, { withFileTypes: true, recursive: true }))
        .filter(file => file.name.endsWith(".tsx"))
        .map(file => path.join(file.parentPath, file.name));

    // Create a TypeScript program
    const program = ts.createProgram(files, {
        target: ts.ScriptTarget.ESNext,
        module: ts.ModuleKind.ESNext,
        jsx: ts.JsxEmit.React,
        moduleResolution: ts.ModuleResolutionKind.NodeNext,
    });

    const typeChecker = program.getTypeChecker();

    // Check if regeneration is needed
    if (!(await shouldRegenerateProps(program, typeChecker))) return;

    console.log("Changes detected in Props interfaces. Regenerating...");

    let output = `// Code generated by props-types-gen. DO NOT EDIT.
package props

import "encoding/json"

`;

    for (const sourceFile of program.getSourceFiles()) {
        if (!sourceFile.fileName.endsWith(".tsx")) continue;

        // Find all interfaces
        ts.forEachChild(sourceFile, node => {
            if (ts.isInterfaceDeclaration(node) && node.name.text.endsWith("Props")) {
                const type = typeChecker.getTypeAtLocation(node);
                const properties = type.getProperties();
                output += generateGoType(node.name.text, properties, typeChecker);
                output += "\n\n";
            }
        });
    }

    // Write output
    const outputDir = "./web/generated";
    await Bun.write(path.join(outputDir, "react_component_props.go"), output);
}

function generateInterfaceHash(sourceFile: ts.SourceFile, typeChecker: ts.TypeChecker): string {
    let hash = "";
    ts.forEachChild(sourceFile, node => {
        if (ts.isInterfaceDeclaration(node) && node.name.text.endsWith("Props")) {
            const type = typeChecker.getTypeAtLocation(node);
            const properties = type.getProperties();
            // Create a string representation of the interface
            hash += `${node.name.text}:{${properties
                .map(prop => {
                    const propType = typeChecker.getTypeOfSymbol(prop);
                    return `${prop.name}:${typeChecker.typeToString(propType)}`;
                })
                .sort()
                .join(",")}}`;
        }
    });
    return hash;
}

async function shouldRegenerateProps(program: ts.Program, typeChecker: ts.TypeChecker): Promise<boolean> {
    const CACHE_FILE = "tmp/props-hash-cache";
    let previousHash = "";

    try {
        previousHash = await Bun.file(CACHE_FILE).text();
    } catch {
        // File doesn't exist, continue
    }

    // Generate new hash from all Props interfaces
    let currentHash = "";
    for (const sourceFile of program.getSourceFiles()) {
        if (!sourceFile.fileName.endsWith(".tsx")) continue;
        currentHash += generateInterfaceHash(sourceFile, typeChecker);
    }

    // Save new hash
    await Bun.write(CACHE_FILE, currentHash);

    // Compare hashes
    return previousHash !== currentHash;
}

function capitalize(str: string): string {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

generateGoTypes().catch(console.error);
