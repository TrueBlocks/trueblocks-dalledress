import { useEffect, useState } from "react";
import { useHotkeys } from "react-hotkeys-hook";

type FetchItemsFn<T> = (curItem: number, perPage: number) => Promise<T[]>;
type FetchCountFn = () => Promise<number>;

export function useKeyboardPaging<T>(
  fetchItems: FetchItemsFn<T>,
  fetchCount: FetchCountFn,
  dependency?: any,
  perPage: number = 20
) {
  const [items, setItems] = useState<T[]>([]);
  const [nItems, setNItems] = useState<number>(0);
  const [curItem, setCurItem] = useState<number>(0);

  useHotkeys("left", (event) => {
    setCurItem((cur) => Math.max(cur - 1, 0));
    event.preventDefault();
  });
  useHotkeys("up", (event) => {
    setCurItem((cur) => Math.max(cur - perPage, 0));
    event.preventDefault();
  });
  useHotkeys("home", (event) => {
    setCurItem(0);
    event.preventDefault();
  });

  useHotkeys("right", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + 1));
    event.preventDefault();
  });
  useHotkeys("down", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem((cur) => Math.min(max, cur + perPage));
    event.preventDefault();
  });
  useHotkeys("end", (event) => {
    var max = Math.max(nItems - perPage, 0);
    setCurItem(max);
    event.preventDefault();
  });

  useEffect(() => {
    const fetch = async () => {
      const fetchedItems = await fetchItems(curItem, perPage);
      console.log("fetchedItems", fetchedItems.length);
      setItems(fetchedItems);
      const cnt = await fetchCount();
      console.log("cnt", cnt);
      setNItems(cnt);
    };
    fetch();
  }, [curItem, perPage, dependency]);

  useEffect(() => {
    setCurItem(0);
  }, [dependency]);

  return { items, nItems, curItem };
}
