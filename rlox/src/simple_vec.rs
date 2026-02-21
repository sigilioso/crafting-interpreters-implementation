use std::{
    alloc::{Layout, alloc, dealloc, handle_alloc_error, realloc},
    ops::Index,
    ptr,
};

pub struct SimpleVec<T> {
    ptr: *mut T,
    capacity: usize,
    count: usize,
}

impl<T> SimpleVec<T> {
    pub fn new() -> Self {
        assert!(
            std::mem::size_of::<T>() != 0,
            "ZST (Zero sized types) are not supported"
        );
        Self {
            ptr: ptr::null_mut(),
            capacity: 0,
            count: 0,
        }
    }

    pub fn with_initial_capacity(capacity: usize) -> Self {
        let mut v = Self::new();
        let layout = Layout::array::<T>(capacity).expect("failure reserving memory");
        let ptr = unsafe { alloc(layout) } as *mut T;
        if ptr.is_null() {
            handle_alloc_error(layout);
        }
        Self {
            ptr,
            capacity,
            count: 0,
        }
    }

    pub fn count(&self) -> usize {
        self.count
    }

    pub fn push(&mut self, item: T) {
        if self.capacity < self.count + 1 {
            self.realloc();
        }
        unsafe {
            ptr::write(self.ptr.add(self.count), item);
        }
        self.count += 1
    }

    pub fn pop(&mut self) -> Option<T> {
        if self.count == 0 {
            None
        } else {
            unsafe {
                self.count -= 1;
                Some(ptr::read(self.ptr.add(self.count)))
            }
        }
    }

    fn grow_capacity(&self) -> usize {
        if self.capacity == 0 {
            8
        } else {
            self.capacity * 2
        }
    }

    fn realloc(&mut self) {
        let capacity = self.grow_capacity();
        let layout = Layout::array::<T>(capacity).expect("failure reserving memory");
        let new_ptr = if self.capacity == 0 {
            unsafe { alloc(layout) }
        } else {
            let old_layout = Layout::array::<T>(self.capacity).expect("failure reserving memory");
            unsafe { realloc(self.ptr as *mut u8, old_layout, layout.size()) }
        };
        if new_ptr.is_null() {
            handle_alloc_error(layout);
        }
        self.ptr = new_ptr as *mut T;
        self.capacity = capacity;
    }

    fn dealloc(&mut self) {
        if !self.ptr.is_null() && self.capacity > 0 {
            while self.pop().is_some() {} // Drop all underlying elements
            let layout = Layout::array::<T>(self.capacity).expect("failure freeing memory");
            unsafe {
                dealloc(self.ptr as *mut u8, layout);
            }
        }
    }
}

impl<T> Drop for SimpleVec<T> {
    fn drop(&mut self) {
        self.dealloc();
    }
}

impl<T> Index<usize> for SimpleVec<T> {
    type Output = T;

    fn index(&self, index: usize) -> &Self::Output {
        unsafe { &*self.ptr.add(index) }
    }
}

#[cfg(test)]
mod tests {
    use crate::value::Value;

    use super::*;

    #[test]
    fn test_simple_vec() {
        let mut v = SimpleVec::<String>::new();
        for i in 1..10 {
            v.push(format!("{i}"));
        }

        assert_eq!(v.capacity, 16);

        for i in 0..9 {
            assert_eq!(v[i].clone(), format!("{}", i + 1));
        }

        let mut v = SimpleVec::<i64>::new();
        v.push(10);
        assert_eq!(v[0], 10);

        let mut v = SimpleVec::<Value>::new();
        v.push(Value::Number(1.2));
        assert!(matches!(v.pop(), Some(Value::Number(n)) if n == 1.2));
        assert_eq!(v.count(), 0);

        let mut v = SimpleVec::<i64>::with_initial_capacity(32);
        v.push(10);
        assert_eq!(v[0], 10);
        assert_eq!(v.capacity, 32);
    }
}
