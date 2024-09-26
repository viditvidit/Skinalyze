/* creating database*/
CREATE DATABASE Skinalyze;
USE Skinalyze;

/* creating tables*/
CREATE TABLE Concern (Concern_ID int primary key, Concern nvarchar(100));
CREATE TABLE Skin_Type (Skin_Type_ID int primary key, Skin_Type nvarchar(100));
CREATE TABLE Brand (Brand_ID int primary key, Brand nvarchar(50));
CREATE TABLE Product_Type (Product_Type_ID int primary key, Product_Type nvarchar(50));
CREATE TABLE Key_Ingredients (Key_Ingredients_ID int primary key, Key_Ingredients nvarchar(50));
CREATE TABLE Products (Product_ID int primary key, Product_Name nvarchar(200),
                       All_Ingredients nvarchar(10000), Concern_ID int,
                       Skin_Type_ID int, Brand_ID int,
                       Product_Type_ID int, Key_Ingredients_ID int,
                       FOREIGN KEY(Concern_ID) REFERENCES Concern(Concern_ID),
                       FOREIGN KEY(Skin_Type_ID) REFERENCES Skin_Type(Skin_Type_ID),
                       FOREIGN KEY(Brand_ID) REFERENCES Brand(Brand_ID),
                       FOREIGN KEY(Product_Type_ID) REFERENCES Product_Type(Product_Type_ID),
                       FOREIGN KEY(Key_Ingredients_ID) REFERENCES Key_Ingredients(Key_Ingredients_ID));
select * from Products;