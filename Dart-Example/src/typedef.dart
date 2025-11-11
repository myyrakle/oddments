typedef void intint(int l, int r);

void addprint(int l, int r)
{
    print(l+r);
}

void mulprint(int l, int r)
{
    print(l*r);
}

main()
{
    intint printer;

    printer = addprint;
    printer(5, 6);

    printer = mulprint;
    printer(5, 6);
}
