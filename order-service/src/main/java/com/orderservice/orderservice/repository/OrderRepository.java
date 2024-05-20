package com.orderservice.orderservice.repository;

import com.orderservice.orderservice.model.Order;
import org.springframework.data.jpa.repository.JpaRepository;

import java.beans.JavaBean;

public interface OrderRepository extends JpaRepository<Order, Long> {
}
