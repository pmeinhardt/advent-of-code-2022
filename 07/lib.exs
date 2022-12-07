defmodule FS do
  defstruct [:pwd, :tree]

  def new do
    %__MODULE__{pwd: "/", tree: %{}}
  end

  def cd(%{pwd: pwd} = fs, path) do
    %{fs | pwd: Path.expand(path, pwd)}
  end

  def mkdir(%{pwd: pwd, tree: tree} = fs, name) do
    %{fs | tree: insert(tree, pwd, name, %{})}
  end

  def touch(%{pwd: pwd, tree: tree} = fs, name, size) do
    %{fs | tree: insert(tree, pwd, name, size)}
  end

  defp insert(%{} = tree, path, key, value) when is_binary(path) do
    components = path |> Path.split() |> Enum.filter(&(&1 != "/"))
    insert(tree, components, key, value)
  end

  defp insert(%{} = tree, [], key, value) do
    Map.put(tree, key, value)
  end

  defp insert(%{} = tree, [name | rest], key, value) do
    Map.put(tree, name, insert(tree[name], rest, key, value))
  end
end

defmodule Reader do
  def parse(input) do
    Enum.reduce(input, FS.new, fn line, fs ->
      case String.trim(line) do
        "$ " <> cmd ->
          case cmd do
            "cd " <> path -> FS.cd(fs, path)
            "ls" -> fs # no effect
          end

        "dir " <> name -> FS.mkdir(fs, name)

        file_info ->
          [size, name] = String.split(file_info, " ")
          FS.touch(fs, name, String.to_integer(size))
      end
    end)
  end
end

defmodule Stats do
  def dirstats(%{} = tree, path \\ "/") do
    {size, stats} = Enum.reduce(tree, {0, %{}}, fn
      {_, file_size}, {sum, st} when is_number(file_size) ->
        {sum + file_size, st}
      {name, sub_tree}, {sum, st} ->
        child_path = Path.join(path, name)
        child_stats = dirstats(sub_tree, child_path)
        {sum + child_stats[child_path], Map.merge(st, child_stats)}
    end)

    Map.put(stats, path, size)
  end
end
