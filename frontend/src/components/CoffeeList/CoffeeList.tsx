import React from "react";
import useSWR from "swr";

import "./style.css";

interface Product {
  readonly id: number;
  name: string;
  price: number;
  sku: string;
  description?: string;
}

function useProductList() {
  const { data, error } = useSWR<Product[]>(
    process.env.REACT_APP_API_ENDPOINT + "/products"
  );

  return {
    products: data,
    isLoading: !error && !data,
    error,
  };
}

const CoffeeList = () => {
  const { products, isLoading } = useProductList();
  const renderProductList = () =>
    products?.map((product) => (
      <tr key={product.id}>
        <td>{product.id}</td>
        <td>{product.name}</td>
        <td>{product.price}</td>
        <td>{product.sku}</td>
        <td>{product.description}</td>
      </tr>
    ));

  return (
    <div>
      <h3>Menu</h3>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Price</th>
            <th>SKU</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {isLoading ? (
            <tr>
              <td>Loading...</td>
            </tr>
          ) : (
            renderProductList()
          )}
        </tbody>
      </table>
    </div>
  );
};

export default CoffeeList;
