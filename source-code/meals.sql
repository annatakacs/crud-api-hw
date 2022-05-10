CREATE USER meals_user WITH PASSWORD 'crud_password';
GRANT ALL PRIVILEGES ON DATABASE meals TO meals_user;
DROP SCHEMA IF EXISTS meals_schema CASCADE;
CREATE SCHEMA meals_schema AUTHORIZATION meals_user;
ALTER SCHEMA meals_schema OWNER TO meals_user;
ALTER ROLE meals_user SET search_path to 'meals_schema', 'public';

CREATE TABLE meals_schema.meals_table (
  id SERIAL PRIMARY KEY,
  name varchar(30) NOT NULL,
  price numeric NOT NULL,
  ingredients varchar(50) NOT NULL,
  spicy boolean NOT NULL,
  vegan boolean NOT NULL,
  gluten_free boolean  NOT NULL,
  description varchar(50) NOT NULL,
  kcal int NOT NULL
) ; 

INSERT INTO meals_schema.meals_table (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) 
VALUES ('Hamburger', 8.45, 'meat, onions, bun, ketchup, mayo', false, false, false, 'Delicious classic burger', 650)
RETURNING id;

INSERT INTO meals_schema.meals_table (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) 
VALUES ('Vegan burger', 9.45, 'mushrooms, onions, bun, ketchup, mayo', false, true, false, 'Delicious vegan burger', 450)
RETURNING id;

INSERT INTO meals_schema.meals_table (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) 
VALUES ('Chili burger', 10.45, 'meat, onions, bun, ketchup, mayo, jalapeno', true, false, false, 'Spicy mexican burger', 700)
RETURNING id;

INSERT INTO meals_schema.meals_table (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) 
VALUES ('Fries', 4.75, 'potato, ketchup', false, true, true, 'French fries the way we love it', 350)
RETURNING id;

INSERT INTO meals_schema.meals_table (name, price, ingredients, spicy, vegan, gluten_free, description, kcal) 
VALUES ('Greek salad', 6.15, 'tomato, lettuce, cucumber, olives, feta cheese', false, false, true, 'Traditional Greek salad', 250)
RETURNING id;