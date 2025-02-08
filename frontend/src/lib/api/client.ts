import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// API response types based on swagger.json
export interface Node {
  id: string;
  name: string;
  nodeData: string;
  createdAt: string;
  updatedAt: string;
  formularNodes: FormularNode[];
}

export interface Formular {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  nodes: FormularNode[];
  calculationFormulars: CalculationFormular[];
}

export interface Calculation {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
  formulars: CalculationFormular[];
}

export interface FormularNode {
  id: string;
  formularId: string;
  nodeId: string;
  nextId: string | null;
  createdAt: string;
  updatedAt: string;
  node: Node;
  formular: Formular;
  next: FormularNode | null;
}

export interface CalculationFormular {
  id: string;
  calculationId: string;
  formularId: string;
  nextId: string | null;
  createdAt: string;
  updatedAt: string;
  calculation: Calculation;
  formular: Formular;
  next: CalculationFormular | null;
}

// API endpoints
export const api = {
  nodes: {
    list: () => apiClient.get<Node[]>('/nodes'),
    get: (id: string) => apiClient.get<Node>(`/nodes/${id}`),
    create: (data: { name: string; nodeData: string }) => 
      apiClient.post<Node>('/nodes', data),
    update: (id: string, data: { name?: string; nodeData?: string }) =>
      apiClient.put<Node>(`/nodes/${id}`, data),
    delete: (id: string) => apiClient.delete(`/nodes/${id}`),
  },
  formulars: {
    list: () => apiClient.get<Formular[]>('/formulars'),
    get: (id: string) => apiClient.get<Formular>(`/formulars/${id}`),
    create: (data: { name: string }) => 
      apiClient.post<Formular>('/formulars', data),
    update: (id: string, data: { name: string }) =>
      apiClient.put<Formular>(`/formulars/${id}`, data),
    delete: (id: string) => apiClient.delete(`/formulars/${id}`),
    getNodes: (id: string) => 
      apiClient.get<FormularNode[]>(`/formulars/${id}/nodes`),
    addNode: (id: string, data: { nodeId: string; nextId?: string }) =>
      apiClient.post<FormularNode>(`/formulars/${id}/nodes`, data),
    removeNode: (formularId: string, nodeId: string) =>
      apiClient.delete(`/formulars/${formularId}/nodes/${nodeId}`),
    reorderNodes: (id: string, data: { nodeOrder: string[] }) =>
      apiClient.put(`/formulars/${id}/nodes/reorder`, data),
  },
  calculations: {
    list: () => apiClient.get<Calculation[]>('/calculations'),
    get: (id: string) => apiClient.get<Calculation>(`/calculations/${id}`),
    create: (data: { name: string }) => 
      apiClient.post<Calculation>('/calculations', data),
    update: (id: string, data: { name: string }) =>
      apiClient.put<Calculation>(`/calculations/${id}`, data),
    delete: (id: string) => apiClient.delete(`/calculations/${id}`),
    getFormulars: (id: string) =>
      apiClient.get<CalculationFormular[]>(`/calculations/${id}/formulars`),
    addFormular: (id: string, data: { formularId: string; nextId?: string }) =>
      apiClient.post<CalculationFormular>(`/calculations/${id}/formulars`, data),
    removeFormular: (calculationId: string, formularId: string) =>
      apiClient.delete(`/calculations/${calculationId}/formulars/${formularId}`),
    reorderFormulars: (id: string, data: { formularOrder: string[] }) =>
      apiClient.put(`/calculations/${id}/formulars/reorder`, data),
  },
};
