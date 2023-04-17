select o.id from "orders" o join order_batches ob on o.batch_id = o.id where ob.valid_until > now();

select o.id, s.status from "orders" o
     left join order_batches ob on ob.id = o.batch_id
     join statuses s on o.status_id = s.id
     where o.status_id = '3'
     and o.pharmacy_id = '1'
     and (o.batch_id = '' or ob.owned_by = '2' or o.id not in
    (select o.id from "orders" o
        join order_batches ob on ob.id = o.batch_id
        where ob.valid_until > now()));

select o.id, ob.valid_until, now() from "orders" o
                     join order_batches ob on ob.id = o.batch_id
where ob.valid_until > now();

select o.id, s.status from "orders" o
                               left join order_batches ob on ob.id = o.batch_id
                               join statuses s on o.status_id = s.id
where o.status_id = '3'
  and o.pharmacy_id = '1'
  and (o.batch_id = '' or ob.owned_by = '2' or o.id not in
                                              (select o.id from "orders" o
                                                                    join order_batches ob on ob.id = o.batch_id
                                               where ob.valid_until > now()));