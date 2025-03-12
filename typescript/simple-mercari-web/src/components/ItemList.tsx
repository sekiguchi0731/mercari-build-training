import { useEffect, useState } from 'react';
import { Item, fetchItems } from '~/api';

const SERVER_URL = import.meta.env.VITE_BACKEND_URL || 'http://127.0.0.1:9001';

interface Prop {
  reload: boolean;
  onLoadCompleted: () => void;
}

export const ItemList = ({ reload, onLoadCompleted }: Prop) => {
  const [items, setItems] = useState<Item[]>([]);
  useEffect(() => {
    const fetchData = () => {
      fetchItems()
        .then((data) => {
          console.debug('GET success:', data);
          setItems(data.items);
          onLoadCompleted();
        })
        .catch((error) => {
          console.error('GET error:', error);
        });
    };

    if (reload) {
      fetchData();
    }
  }, [reload, onLoadCompleted]);

  return (
    <div className='ItemField'>
      {items?.map((item) => {
        return (
          <div key={item.id} className="ItemList">
            {/* Show item images */}
            <img
              src={`${SERVER_URL}/${item.image_name}`}
              alt={item.name}
            />
            <p>
              <span>ID: {item.id}</span>
              <br />
              <span>Name: {item.name}</span>
              <br />
              <span>Category: {item.category}</span>
            </p>
          </div>
        );
      })}
    </div>
  );
};
